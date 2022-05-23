package main

import (
	"GoSearchEngine/fenci"
	"GoSearchEngine/storage"
	"fmt"
	"strings"
)

func main() {
	fenci.ReadDict()
	// 此句进行分词，分词结果和文档存储在 database 文件夹中
	// 若需要重新分词请先清空 database 文件夹
	// utils.WukongFenCi()

	// 循环接收输入，使用 Ctrl+C 中断即可
	searchDemo()
}

func searchDemo() {
	key := ""
	for {
		fmt.Print("请输入关键词：")
		fmt.Scanf("%s\n", &key)

		// 获取关键词索引的 ID 列表
		idByteList, _ := storage.DictDB.Get([]byte(key))
		idList := string(idByteList)

		// 获取文档ID列表
		// 由于存储时使用逗号作为分隔。所以读取时使用逗号分割
		result := strings.Split(idList, ",")
		fmt.Println(len(result) - 1)
		fmt.Println(idList)

		// 输出所有文档
		for _, id := range result {
			doc, _ := storage.DocDB.Get([]byte(id))
			fmt.Println(id, string(doc))
		}
	}
}
