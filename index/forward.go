package index

import (
	"GoSearchEngine/storage"
	json2 "encoding/json"
	"log"
	"strconv"
)

type ForwardStore struct {
	Keywords []string `json:"keywords"`
	Times    []int    `json:"times"`
}

// AddWordsToForwardIndex 为词列表和次数列表添加正排索引
func AddWordsToForwardIndex(id int, words []string, times []int) {
	f := ForwardStore{
		words,
		times,
	}
	f.save(id)
}

// 将关键词列表和次数列表存入数据库
func (store ForwardStore) save(id int) {
	json, _ := json2.Marshal(store)
	err := storage.ForwardIndex.Set([]byte(strconv.Itoa(id)), json)
	if err != nil {
		log.Fatalln(id, json, " SET 失败", err)
	}
}

// GetIdWords 获取 ID 对应的文档的关键词列表
func GetIdWords(id int) (ForwardStore, bool) {
	wordsListJson, e := storage.ForwardIndex.Get([]byte(strconv.Itoa(id)))
	if !e {
		return ForwardStore{}, false
	}

	wordsList := ForwardStore{}
	err := json2.Unmarshal(wordsListJson, &wordsList)
	if err != nil {
		log.Fatalln("错误的 json ", err)
	}
	return wordsList, true
}

// GetForwardIndex 获取指定 ID 文档的正排索引
func GetForwardIndex(id int) (ForwardStore, bool) {
	return GetIdWords(id)
}
