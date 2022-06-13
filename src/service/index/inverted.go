package index

import (
	"bytes"
	"fmt"
	"go-search-engine/src/service/storage"
	"log"
	"strconv"
	"strings"
	"sync"
)

type InvertedIndexGen struct {
	cache *sync.Map
}

func (i *InvertedIndexGen) AddWordsToInvertedIndex(words []string, id int) {
	for _, keyWord := range words {
		// 将所有分词保存到倒排索引中
		if len(keyWord) == 0 {
			continue
		}
		i.SaveToCache(keyWord, id)
	}
}

func (i *InvertedIndexGen) SaveToCache(word string, id int) {
	v, _ := i.cache.LoadOrStore(word, &bytes.Buffer{})
	v.(*bytes.Buffer).WriteString(fmt.Sprintf("%d,", id))
}

func (i *InvertedIndexGen) Flush() {
	i.cache.Range(func(key, val any) bool {
		err := storage.InvertedIndex.Set([]byte(key.(string)), val.(*bytes.Buffer).Bytes())
		if err != nil {
			log.Println("ivertered index flush set err: ", err)
			return false
		}
		return true
	})
}

// GetWordIds 获取关键词索引的 ID 列表
func GetWordIds(word string) ([]int, bool) {
	idByteList, e := storage.InvertedIndex.Get([]byte(word))
	if !e {
		return nil, false
	}
	idList := string(idByteList)

	// 由于存储时使用逗号作为分隔。所以读取时使用逗号分割
	idStr := strings.Split(idList, ",")
	// 这里减一是因为存入数据库时以 , 结尾，idStr 的末尾元素为空串
	ids := make([]int, len(idStr)-1)
	for i := 0; i < len(idStr)-1; i++ {
		ids[i], _ = strconv.Atoi(idStr[i])
	}
	return ids, true
}

func GetInvertedIndex(word string) ([]int, bool) {
	return GetWordIds(word)
}
