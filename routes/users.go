package routes

import (
	"github.com/gin-gonic/gin"
	. "github.com/gin2/controllers"
	middleware2 "github.com/gin2/middleware"
	"github.com/gin2/pkg/setting"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/gin2/docs"
)

type UserRoute struct {
	userController *UserController
	middleware     *middleware2.Jwt
	auth           *Auth
}

func (userRoute *UserRoute) InitRoute() *gin.Engine {

	app := gin.Default()
	gin.SetMode(setting.RunMode)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/auth",userRoute.auth.Auth)
	group := app.Group("/users")
	group.Use(userRoute.middleware.JwtTokenValidate)
	group.GET("/", userRoute.userController.Index)
	group.POST("/", userRoute.userController.Store)
	group.GET("/:id", userRoute.userController.Show)
	return app
}

func NewUserRoute() *UserRoute {

	route := &UserRoute{}
	route.userController = NewUserController()
	route.middleware = middleware2.NewJwt()
	route.auth = NewAuth()
	return route

}
