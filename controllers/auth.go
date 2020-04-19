package controllers

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin2/models"
	error2 "github.com/gin2/pkg/error"
	"github.com/gin2/pkg/logging"
	"github.com/gin2/pkg/util"
)

type Auth struct {
	UserName string `valid:"Required;MaxSize(32)"`
	Password string `valid:"Required;MaxSize(32)"`
}

func NewAuth() *Auth {

	return &Auth{}
}
func (userController *Auth) Auth(ctx *gin.Context) {

	userName := ctx.Query("user_name")
	password := ctx.Query("password")

	valid := validation.Validation{}
	auth := Auth{UserName: userName, Password: password}
	ok, _ := valid.Valid(&auth)

	data := make(map[string]interface{})
	code := error2.INVALID_PARAMS
	if ok {
		userModel := models.NewUserModel()
		code = error2.ERROR_AUTH
		checkLogin := userModel.CheckLogin(userName, password)
		if checkLogin {
			code = error2.ERROR_AUTH_TOKEN
			token, err := util.GenerateToken(userName, password)
			if err == nil {
				code = error2.SUCCESS
				data["token"] = token

			}
		}
	}

	for _, err := range valid.Errors {
		str := fmt.Sprintf("error key:%s, err message:%s", err.Key, err.Message)
		logging.Info(str)
	}

	ctx.JSON(200, gin.H{
		"code":    code,
		"message": error2.GetErrorMsg(code),
		"data":    data,
	})

}
