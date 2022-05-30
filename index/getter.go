package index

import (
	"GoSearchEngine/storage"
	"GoSearchEngine/utils"
	json2 "encoding/json"
	"log"
	"strconv"
	"strings"
)

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

// GetIdWords 获取 ID 对应的文档的关键词列表
func GetIdWords(id int) ([]string, bool) {
	wordsListJson, e := storage.ForwardIndex.Get([]byte(strconv.Itoa(id)))
	if !e {
		return nil, false
	}

	wordsList := make([]string, 10)
	err := json2.Unmarshal(wordsListJson, &wordsList)
	if err != nil {
		log.Fatalln("错误的 json ", err)
	}
	return wordsList, true
}

func ReadWukong() <-chan utils.CsvRow {
	return utils.ReadCsv(utils.GetPath("/dataset/wukong.csv"), 2, true)
}

func GetDocument(id int) ([]byte, bool) {
	key := []byte(strconv.Itoa(id))
	return storage.DocDB.Get(key)
}
