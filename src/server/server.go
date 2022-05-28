package server

import (
	"go-search-engine/src/config"
	router "go-search-engine/src/router"

	"github.com/gin-gonic/gin"
)

func Run(httpServer *gin.Engine) {

	//设置session
	//gob.Register(global.TMember{})
	//httpServer.Use(global.GetSession())

	// 注册路由
	router.RegisterRouter(httpServer)

	serverError := httpServer.Run(config.GetServerConfig().HTTP_HOST + ":" + config.GetServerConfig().HTTP_PORT)

	if serverError != nil {
		panic("server error !" + serverError.Error())
	}

}
