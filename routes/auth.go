package routes

import (
	"github.com/gin-gonic/gin"
	. "github.com/gin2/controllers"
)

type AuthRoute struct {
	auth *Auth
}

func (authRoute *AuthRoute) LoadAuthRoute(e *gin.Engine) {

	e.GET("/auth", authRoute.auth.Auth)

}

func NewAuthRoute() *AuthRoute {

	route := &AuthRoute{}
	route.auth = NewAuth()
	return route

}
