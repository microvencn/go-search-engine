package search

import (
	"github.com/gin-gonic/gin"
	global "go-search-engine/src/global"
	"go-search-engine/src/service/index"
	"go-search-engine/src/service/related_search"
	"go-search-engine/src/service/searcher"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	filter := strings.Split(c.Query("filter"), ",")
	if len(filter) == 0 {
		filter = nil
	} else {
		for i := 0; i < len(filter); i++ {
			filter[i] = strings.Replace(filter[i], " ", "", -1)
			filter[i] = strings.Replace(filter[i], "\n", "", -1)
		}
	}
	total, data := searcher.SimpleWithFilter(query, offset, pageSize, filter)
	c.JSON(http.StatusOK, global.SearchResponse{Total: total, Data: data})
	return
}

func AutoComplete(c *gin.Context) {
	prefix := c.Query("prefix")
	if len(prefix) == 0 {
		c.JSON(http.StatusOK, global.AutoCompleteResponse{Data: make([]string, 0)})
		return
	}
	wordList := index.TrieTree.Search(prefix, 7)
	if wordList == nil {
		c.JSON(http.StatusOK, global.AutoCompleteResponse{Data: make([]string, 0)})
		return
	}
	result := make([]string, 0, len(*wordList))
	if wordList != nil {
		sort.Sort(wordList)
		for i := 0; i < len(*wordList); i++ {
			result = append(result, (*wordList)[i].Text)
		}
	}
	c.JSON(http.StatusOK, global.AutoCompleteResponse{Data: result})
}

func RelatedWords(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusOK, global.ResponseMeta{})
		return
	}
	data := related_search.GetQueryRelatedWords(query)
	c.JSON(http.StatusOK, global.RelatedSearchResponse{Data: data})
}
