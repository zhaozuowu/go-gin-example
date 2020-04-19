package controllers

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin2/models"
	error2 "github.com/gin2/pkg/error"
	"github.com/gin2/pkg/setting"
	"github.com/gin2/request"
	"github.com/gin2/service"
	"github.com/unknwon/com"
	"log"
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
	data := make(map[string]interface{})
	var userRequest request.UserListRequst
	if err := ctx.ShouldBindQuery(&userRequest); err != nil {
		code := error2.INVALID_PARAMS
		ctx.JSON(code, gin.H{
			"code":  code,
			"data":  "",
			"error": err.Error(),
			"msg":   error2.GetErrorMsg(code),
		})
		return
	}
	page := com.StrTo(ctx.DefaultQuery("page", "1")).MustInt()
	requstParams["name"] = userRequest.Name
	requstParams["email"] = userRequest.Email

	userModel := models.NewUserModel()
	offset := 0
	offset = (page - 1) * setting.PageSize
	if offset < 0 {
		offset = 0
	}

	userList, err := userModel.GetUserList(offset, setting.PageSize, requstParams)
	if err != nil {
		code := error2.ERROR
		ctx.JSON(code, gin.H{
			"code":  code,
			"data":  "",
			"error": err.Error(),
			"msg":   error2.GetErrorMsg(code),
		})
		return
	}

	total, err := userModel.GetUserTotalNum(requstParams)
	if err != nil {
		code := error2.ERROR
		ctx.JSON(code, gin.H{
			"code":  code,
			"data":  "",
			"error": err.Error(),
			"msg":   error2.GetErrorMsg(code),
		})
		return
	}

	code := error2.SUCCESS
	data["userList"] = userList
	data["total"] = total
	ctx.JSON(code, gin.H{
		"code": code,
		"data": data,
		"msg":  error2.GetErrorMsg(code),
	})

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

	id := com.StrTo(ctx.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Required(id, "id").Message("id参数必传")
	valid.Min(id, 1, "id").Message("id必须大于0")
	code := error2.SUCCESS
	result := map[string]interface{}{"code": code, "message": error2.GetErrorMsg(code), "data": ""}
	if valid.HasErrors() {
		code = error2.INVALID_PARAMS
		for _, err := range valid.Errors {
			//log.Fatalf("error key:%s, err message:%s\n", err.Key, err.Message)
			log.Printf("error key:%s, err message:%s\n", err.Key, err.Message)
		}
		result["code"] = code
		result["message"] = error2.GetErrorMsg(code)
		ctx.JSON(200, result)
		return
	}
	userModel := models.NewUserModel()
	userInfo, err := userModel.GetUserInfoById(id)
	result["data"] = userInfo
	result["err"] = err
	result["debug"] = "debug2"
	ctx.JSON(200, result)
	return

}

func NewUserController() *UserController {

	user := &UserController{}
	user.userService = service.NewUserService()
	return user
}
