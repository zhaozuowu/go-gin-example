package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin2/pkg/setting"
	"github.com/gin2/routes"
	"log"
	"syscall"
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/zhaozuowu/go-gin-example
// @license.name MIT
// @license.url https://github.com/zhaozuowu/go-gin-example/blob/master/LICENSE
func main() {

	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.HttpPort)

	userRoute := routes.NewUserRoute()
	authRoute := routes.NewAuthRoute()
	routes.IncludeRoute(userRoute.LoadUser,authRoute.LoadAuthRoute)
	server := endless.NewServer(endPoint, routes.InitRoute())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
