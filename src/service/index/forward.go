package index

import (
	json2 "encoding/json"
	"go-search-engine/src/service/storage"
	"log"
	"strconv"
)

// AddWordsToForwardIndex 为词列表和次数列表添加正排索引
func AddWordsToForwardIndex(id int, words []string, times []int, topKWords []string, topKWeights []float64) {
	f := storage.ForwardStore{
		Keywords:    words,
		Times:       times,
		TopKWords:   topKWords,
		TopKWeights: topKWeights,
	}
	f.Save(id)
}

// GetIdWords 获取 ID 对应的文档的关键词列表
func GetIdWords(id int) (storage.ForwardStore, bool) {
	wordsListJson, e := storage.ForwardIndex.Get([]byte(strconv.Itoa(id)))
	if !e {
		return storage.ForwardStore{}, false
	}

	wordsList := storage.ForwardStore{}
	err := json2.Unmarshal(wordsListJson, &wordsList)
	if err != nil {
		log.Fatalln("错误的 json ", err)
	}
	return wordsList, true
}

// GetForwardIndex 获取指定 ID 文档的正排索引
func GetForwardIndex(id int) (storage.ForwardStore, bool) {
	return GetIdWords(id)
}
