package search

import (
	"github.com/gin-gonic/gin"
	global "go-search-engine/src/global"
	"go-search-engine/src/service/searcher"
	"net/http"
	"strconv"
)

func SimpleSearch(c *gin.Context) {
	query := c.Query("query")
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err2 := strconv.Atoi(c.Query("page_size"))
	if err != nil || err2 != nil {
		c.JSON(http.StatusOK, global.UnknownError)
		return
	}
	offset := (page - 1) * pageSize
	total, data := searcher.SimpleWithFilter(query, offset, pageSize, make([]string, 0))
	c.JSON(http.StatusOK, global.SearchResponse{Total: total, Data: data})
	return
}

func TrieSearch(c *gin.Context) {
	query := c.Query("query")
	trieList := searcher.SimpleTrieSearch(query)
	c.JSON(http.StatusOK, trieList)
	return
}
