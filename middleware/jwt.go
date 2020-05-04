package middleware

import (
	"github.com/gin-gonic/gin"
	error2 "github.com/gin2/pkg/error"
	"github.com/gin2/pkg/util"
	"net/http"
	"time"
)

type Jwt struct {
}

func (jwt *Jwt) JwtTokenValidate(ctx *gin.Context) {

	code := error2.SUCCESS
	var data interface{}
	token := ctx.Query("token")
	if token == "" {
		code = error2.TOKEN_IS_EMPTY
	} else {
		claims, err := util.ParseToken(token)
		if err != nil {
			code = error2.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = error2.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		}
	}
	if code != error2.SUCCESS {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    code,
			"message": error2.GetErrorMsg(code),
			"data":    data,
		})
		ctx.Abort()
		return
	}

	ctx.Next()

}
func NewJwt() *Jwt {

	return &Jwt{}
}
