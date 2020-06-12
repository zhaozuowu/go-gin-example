package controllers

import (
	"fmt"
	"github.com/gin2/pkg/redis"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin2/models"
	"github.com/gin2/pkg/app"
	error2 "github.com/gin2/pkg/error"
	"github.com/gin2/pkg/logging"
	"github.com/gin2/request"
	"github.com/gin2/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/unknwon/com"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

type Result struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"请求成功"`
	Data    interface{} `json:"data"`
}

// @Summary 用户列表
// @Id 1
// @Tags 用户中心
// @version 1.0
// @Accept application/json
// @Param name body string  true "Name"
// @Param email body string  true "Email"
// @Param page query string false "Page"
// @Success 200 object controllers.Result
// @Router /users/ [get]
func (userController *UserController) Index(ctx *gin.Context) {

	requstParams := make(map[string]interface{})
	//var userRequest request.UserListRequst
	appG := app.Gin{ctx}
	/*
	if err := ctx.ShouldBindQuery(&userRequest); err != nil {
		code := error2.INVALID_PARAMS
		appG.Response(http.StatusOK, code, nil)
		return
	}

	valid := validation.Validation{}
	valid.MinSize(userRequest.Name, 5, "name").Message("用户名称不能少于10个字符")
	valid.MinSize(userRequest.Email, 5, "email").Message("邮箱长度不能少于10个字符")

	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.Response(http.StatusOK, error2.INVALID_PARAMS, nil)
		return
	}*/

	page := com.StrTo(ctx.DefaultQuery("page", "1")).MustInt()
	//requstParams["name"] = userRequest.Name
	//requstParams["email"] = userRequest.Email
	userList, err := userController.userService.GetUserList(requstParams, page)

	if err != nil {
		logging.Error("GetUserList error",err)
		appG.Response(http.StatusOK,error2.ERROR_GET_USERLIST_FAIL,nil)
		return
	}

	appG.Response(http.StatusOK,error2.SUCCESS,userList)

}

// @Summary 用户注册接口
// @Id 1
// @Tags 用户中心
// @version 1.0
// @Accept application/json
// @Param name body string  true "Name"
// @Param email body string  true "Email"
// @Success 200 object controllers.Result
// @Router /users/ [post]
func (userController *UserController) Store(ctx *gin.Context) {

	var userRequst request.UserListRequst

	if err := ctx.ShouldBindJSON(&userRequst); err != nil {
		code := error2.INVALID_PARAMS

		ctx.JSON(code, gin.H{
			"code":    code,
			"message": error2.GetErrorMsg(code),
			"error":   err.Error(),
			"data":    "",
		})
		return
	}

	name := userRequst.Name
	email := userRequst.Email
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 32, "name").Message("名称最长不能超过32个字符")
	valid.Required(email, "email").Message("邮箱不能为空")
	valid.MaxSize(email, 32, "email").Message("邮箱最长不能超过32个字符")
	code := error2.INVALID_PARAMS
	if valid.HasErrors() {
		ctx.JSON(200, gin.H{
			"code":    code,
			"message": error2.GetErrorMsg(code),
			//"error":   valid.Errors,
			"data": "",
		})
		return
	}

	userModel := models.NewUserModel()
	if result := userModel.CheckEmailHaveExists(email); result {
		code = error2.EMAIL_HAS_EXISTS
		ctx.JSON(200, gin.H{
			"code":    code,
			"message": error2.GetErrorMsg(code),
			//"error":   valid.Errors,
			"data": "",
		})
		return
	}

	if err := userModel.CreateUser(name, email); err != nil {

		code := error2.REGISTER_USER__FAIL
		ctx.JSON(200, gin.H{
			"code":    code,
			"message": error2.GetErrorMsg(code),
			"data":    "",
			"error":   err,
		})
		return
	}

	code = error2.SUCCESS
	ctx.JSON(200, gin.H{
		"code":    code,
		"message": error2.GetErrorMsg(code),
		"data":    "",
	})
	return

}

// @Summary 获取用户信息
// @Id 1
// @Tags 用户中心
// @version 1.0
// @Accept application/json
// @Param id query int  true "Id"
// @Success 200 object controllers.Result
// @Router /users/:id [get]
func (userController *UserController) Show(ctx *gin.Context) {
	code := error2.SUCCESS
	result := map[string]interface{}{"code": code, "message": error2.GetErrorMsg(code), "data": ""}

	hystrix.ConfigureCommand(ctx.Request.URL.Path+"."+ctx.Request.Method, hystrix.CommandConfig{
		Timeout:               100,
		MaxConcurrentRequests: 600,
		ErrorPercentThreshold: 10,
		SleepWindow:           5000,
		RequestVolumeThreshold:20,
	})

	hystrix.Do(ctx.Request.URL.Path+"."+ctx.Request.Method, func() error {

		id := com.StrTo(ctx.Param("id")).MustInt()
		valid := validation.Validation{}
		valid.Required(id, "id").Message("id参数必传")
		valid.Min(id, 1, "id").Message("id必须大于0")
		if valid.HasErrors() {
			code = error2.INVALID_PARAMS
			for _, err := range valid.Errors {
				//log.Fatalf("error key:%s, err message:%s\n", err.Key, err.Message)
				logging.Error(fmt.Sprintf("error key:%s, err message:%s\n", err.Key, err.Message))
			}
			result["code"] = code
			result["message"] = error2.GetErrorMsg(code)
			ctx.JSON(200, result)
			return nil
		}

		cacheKey := fmt.Sprintf("seckill:v2:user:%d",id)
		cacheResult, err := redis.Get(cacheKey)

		var userInfo models.User
		if err == nil {
			destring := redis.Gzdecode(cacheResult)
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			json.Unmarshal(destring,&userInfo)
			if userInfo.ID >0  {
				result["data"] = userInfo
				result["err"] = err
				result["debug"] = "debug2"
				ctx.JSON(200, result)
				return nil
			}
		}
		userModel := models.NewUserModel()
		userInfo, err = userModel.GetUserInfoById(id)

		 redis.Set(cacheKey,userInfo,error2.CACHE_EXPIRE_TIME)
		result["data"] = userInfo
		result["err"] = err
		result["debug"] = "debug2"
		ctx.JSON(200, result)
		return nil

	}, func(err error) error {
		logging.Error("get an error,handle it", err,ctx.Request.URL.Path+"."+ctx.Request.Method)
		result["code"] = error2.ERROR
		result["err"] = err
		ctx.JSON(200, result)
		return nil
	})

}

func NewUserController() *UserController {

	user := &UserController{}
	user.userService = service.NewUserService()
	return user
}
