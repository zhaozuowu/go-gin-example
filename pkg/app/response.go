package app

import (
	"github.com/gin-gonic/gin"
	error2 "github.com/gin2/pkg/error"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errorCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code":    errorCode,
		"message": error2.GetErrorMsg(errorCode),
		"data":    data,
	})
}
