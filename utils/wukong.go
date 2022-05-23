package utils

import (
	"GoSearchEngine/avl_struct"
	"GoSearchEngine/fenci"
	"GoSearchEngine/storage"
	"log"
	"strconv"
	"sync"
	"unicode"
)

func ReadWukong() <-chan CsvRow {
	return ReadCsv("./dataset/wukong.csv", 2, true)
}

func WukongFenCi() {
	rows := ReadWukong()
	wg := sync.WaitGroup{}
	// 开启五个线程同时处理分词
	// 这也是为什么 ReadCsv 返回 chan 的原因
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			for csvRow := range rows {
				// 使用文档位于CSV中的 行数-1（忽略表头）作为文档ID
				SaveDocument(csvRow.RowNo, &csvRow.Columns[1])

				// 使用 AVL 对分词后的关键词进行去重
				tree := avl_struct.Init[string]()
				// 存储关键词
				fenci.ExecAndDoSomething(csvRow.Columns[1], func(word string) {
					// 筛选分词结果
					if len(word) == 0 || word == " " {
						return
					}
					runes := []rune(word)
					if unicode.IsPunct(runes[0]) {
						return
					}
					tree.Insert(word)
				})
				// 遍历二叉树得到不重复的关键词列表
				// 并将 <关键词, 文档 ID> 加入数据库
				keyWords := tree.PreOrder()
				for _, keyWord := range keyWords {
					SaveKeyIndex(&keyWord, csvRow.RowNo)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func SaveDocument(id int, doc *string) {
	err := storage.DocDB.Set([]byte(strconv.Itoa(id)), []byte(*doc))
	if err != nil {
		log.Println(doc, "存储失败")
		return
	}
}

func SaveKeyIndex(keyWord *string, id int) {
	bytes := []byte(*keyWord)
	idList := ""

	// 若已存在于数据库中则在其后追加文档ID
	// 目前认为每个文档只会执行一次，所以对结果不进行去重
	// 后面再考虑是否去重
	value, exists := storage.DictDB.Get(bytes)
	if exists {
		idList = string(value)
	}
	// 追加文档ID并写入数据库
	idList += strconv.Itoa(id) + ","
	err := storage.DictDB.Set(bytes, []byte(idList))
	if err != nil {
		log.Println(*keyWord, " SET 失败")
	}
}
