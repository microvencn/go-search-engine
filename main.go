package main

import (
	"github.com/gin-gonic/gin"
	"go-search-engine/src/database"
	global "go-search-engine/src/global"
	"go-search-engine/src/server"
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/index"
	"net/http"
)

func main() {
	fenci.ReadDict()
	index.InitTrie()

	defer func() {
		database.MySqlDb.Close()
		// database.RedisClient.Close()
	}()

	httpServer := gin.Default()
	httpServer.Use(CrosHandler())
	server.Run(httpServer)
}

func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, global.ResponseMeta{Code: http.StatusOK})
		}

		//处理请求
		context.Next()
	}
}
