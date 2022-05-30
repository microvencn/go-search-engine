package index

import (
	"GoSearchEngine/keywords"
	"GoSearchEngine/storage"
	json2 "encoding/json"
	"log"
	"strconv"
)

func SaveDocument(id int, doc *string) {
	err := storage.DocDB.Set([]byte(strconv.Itoa(id)), []byte(*doc))
	if err != nil {
		log.Println(doc, "存储失败")
		return
	}
}

func SaveWordId(keyword string, id int) {
	bytes := []byte(keyword)
	idList := ""

	// 若已存在于数据库中则在其后追加文档ID
	// 目前认为每个文档只会执行一次，所以对结果不进行去重
	// 后面再考虑是否去重
	value, exists := storage.InvertedIndex.Get(bytes)
	if exists {
		idList = string(value)
	} else {
		keywords.AddKeyWords(keyword)
	}

	// 追加文档ID并写入数据库
	idList += strconv.Itoa(id) + ","
	err := storage.InvertedIndex.Set(bytes, []byte(idList))
	if err != nil {
		log.Fatalln(keyword, " SET 失败")
	}
}

// SaveIdWords 将传入的参数与词列表存入 leveldb，若数据库中已存在则直接覆盖
func SaveIdWords(id int, keywords []string) {
	json, _ := json2.Marshal(keywords)
	err := storage.ForwardIndex.Set([]byte(strconv.Itoa(id)), json)
	if err != nil {
		log.Fatalln(id, json, " SET 失败")
	}
}
