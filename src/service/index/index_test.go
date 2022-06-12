package index

import (
	json2 "encoding/json"
	"fmt"
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/storage"
	"go-search-engine/src/service/trie"
	"log"
	"sort"
	"strconv"
	"testing"
)

// 通过运行此测试可以重新生成悟空数据集的索引以及关键词列表
func TestInitWukongIndex(t *testing.T) {
	fenci.ReadDict()
	InitWukongIndex()
	//saveIDF()
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

func TestInvertedIndex(t *testing.T) {
	r, _ := storage.InvertedIndex.Get([]byte("图片"))
	log.Println(r)
}

func TestTrie_Create(t *testing.T) {
	//构造前缀树
	fenci.ReadDict()
	tree := InitTrie()
	trie.WriteTrieFile(tree)
}

func TestTrie_Read(t *testing.T) {
	tree := trie.ReadTrieFile()
	wordList := tree.Search("美", -1)
	if wordList != nil {
		sort.Sort(wordList)
		*wordList = (*wordList)[0:10]
		//fmt.Println(wordList)
		for _, v := range *wordList {
			fmt.Printf("%s %d\n", v.Text, v.Used)
		}
	} else {
		fmt.Println("结果为空")
	}
}
