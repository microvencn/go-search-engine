package main

import (
	"github.com/gin-gonic/gin"
	"go-search-engine/src/database"
	"go-search-engine/src/server"
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/searcher"
	"go-search-engine/src/service/trie"
	"time"
)

func main() {
	fenci.ReadDict()
	//初始化前缀树
	searcher.Tree = trie.ReadTrieFile()
	//每隔一段时间持久化前缀树
	go func() {
		time.AfterFunc(2*time.Hour, func() {
			trie.WriteTrieFile(searcher.Tree)
		})
	}()
	defer func() {
		database.MySqlDb.Close()
		// database.RedisClient.Close()
	}()

	httpServer := gin.Default()
	server.Run(httpServer)
}
