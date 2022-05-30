package index

import (
	"GoSearchEngine/fenci"
	"GoSearchEngine/storage"
	json2 "encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
)

// 通过运行此测试可以重新生成悟空数据集的索引以及关键词列表
func TestInitWukongIndex(t *testing.T) {
	fenci.ReadDict()
	InitWukongIndex()
}

func forwardSearch() {
	key := 1
	for {
		fmt.Print("请输入文档 ID：")
		fmt.Scanf("%d\n", &key)

		// 获取 ID 对应的文档的关键词列表
		wordsListJson, _ := storage.ForwardIndex.Get([]byte(strconv.Itoa(key)))
		wordsList := make([]string, 10)
		err := json2.Unmarshal(wordsListJson, &wordsList)
		if err != nil {
			log.Fatalln("错误的 json ", err)
		}

		// 输出文档
		doc, _ := storage.DocDB.Get([]byte(strconv.Itoa(key)))
		fmt.Println(string(doc))

		// 输出所有关键词
		for _, word := range wordsList {
			fmt.Println(string(word))
		}
	}
}
