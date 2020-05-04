package routes

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/gin2/pkg/setting"
)

type Route func(engine *gin.Engine)

var routes = []Route{}

func IncludeRoute(route ...Route) {

	routes = append(routes, route...)
}

func InitRoute() *gin.Engine {

	app := gin.Default()
	gin.SetMode(setting.RunMode)

	for _, handler := range routes {
		handler(app)
	}
	ginpprof.Wrap(app)
	return app
}
