package main

import (
	"github.com/gin-gonic/gin"
	"go-search-engine/src/database"
	"go-search-engine/src/server"
	"go-search-engine/src/service/fenci"
)

func main() {
	fenci.ReadDict()

	defer func() {
		database.MySqlDb.Close()
		// database.RedisClient.Close()
	}()

	httpServer := gin.Default()
	server.Run(httpServer)
}
