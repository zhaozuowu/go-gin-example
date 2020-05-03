package routes

import (
	"github.com/gin-gonic/gin"
	. "github.com/gin2/controllers"
	_ "github.com/gin2/docs"
	middleware2 "github.com/gin2/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type UserRoute struct {
	userController *UserController
	middleware     *middleware2.Jwt
}

func (userRoute *UserRoute) LoadUser(e *gin.Engine)  {


	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group := e.Group("/users")
	group.Use(userRoute.middleware.JwtTokenValidate)
	group.GET("/", userRoute.userController.Index)
	group.POST("/", userRoute.userController.Store)
	group.GET("/:id", userRoute.userController.Show)

}

func NewUserRoute() *UserRoute {

	route := &UserRoute{}
	route.userController = NewUserController()
	route.middleware = middleware2.NewJwt()
	return route

}
