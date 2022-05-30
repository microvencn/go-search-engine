package index

import (
	"GoSearchEngine/avl_struct"
	"GoSearchEngine/fenci"
	"GoSearchEngine/keywords"
	"strings"
	"sync"
	"unicode"
)

// InitWukongIndex 初始化悟空数据集索引
func InitWukongIndex() {
	keywords.InitKeyWordsFile()
	defer keywords.CloseKeywordsFile()
	rows := ReadWukong()
	wg := sync.WaitGroup{}
	// 开启五个线程同时处理分词
	// 这也是为什么 ReadCsv 返回 chan 的原因
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			for csvRow := range rows {
				doc := strings.ToLower(csvRow.Columns[1])
				words := splitUniqueWords(&doc)
				// 使用文档位于CSV中的 行数-1（忽略表头）作为文档ID
				AddWordsToInvertedIndex(words, csvRow.RowNo)
				AddWordsToForwardIndex(words, csvRow.RowNo)
				SaveDocument(csvRow.RowNo, &csvRow.Columns[1])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// AddDocToInvertedIndex 为一个文档添加倒排索引
func AddDocToInvertedIndex(doc *string, id int) {
	words := splitUniqueWords(doc)
	AddWordsToInvertedIndex(words, id)
}

func AddWordsToInvertedIndex(words []string, id int) {
	for _, keyWord := range words {
		// 将所有分词保存到倒排索引中
		SaveWordId(keyWord, id)
	}
}

// AddDocToForwardIndex 为一个文档添加正排索引
func AddDocToForwardIndex(doc *string, id int) {
	keyWords := splitUniqueWords(doc)
	AddWordsToForwardIndex(keyWords, id)
}

func AddWordsToForwardIndex(words []string, id int) {
	SaveIdWords(id, words)
}

// splitUniqueWords 对文档进行分词，且所有分词结果不重复
func splitUniqueWords(doc *string) []string {
	// 使用 AVL 对分词后的关键词进行去重
	tree := avl_struct.Init[string]()
	fenci.ExecAndDoSomething(doc, func(word string) {
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
	return tree.Inorder()
}
