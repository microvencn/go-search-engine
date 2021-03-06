package index

import (
	"bufio"
	"fmt"
	"go-search-engine/src/service/avl_struct"
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/keywords"
	"go-search-engine/src/service/storage"
	"go-search-engine/src/service/trie"
	"go-search-engine/src/service/utils"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var TrieTree *trie.Trie

var docNum int = 0

// InitWukongIndex 初始化悟空数据集索引
func InitWukongIndex() {
	keywords.InitKeyWordsFile()
	defer keywords.CloseKeywordsFile()
	rows := ReadWukong()
	wg := sync.WaitGroup{}
	iig := InvertedIndexGen{
		cache: &sync.Map{},
	}

	// 记录总的文档数目
	ch := make(chan int)
	defer close(ch)
	go func() {
		for i := range ch {
			docNum += i
		}
	}()

	// 开启五个线程同时处理分词
	// 这也是为什么 ReadCsv 返回 chan 的原因
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func() {
			for csvRow := range rows {
				csvRow.Columns[1] = utils.RemoveSpace(csvRow.Columns[1])
				if len(csvRow.Columns) == 0 {
					continue
				}
				doc := strings.ToLower(csvRow.Columns[1])
				ch <- 1
				words, _ := splitUniqueWords(&doc)
				// 使用文档位于CSV中的 行数-1（忽略表头）作为文档ID
				iig.AddWordsToInvertedIndex(words, csvRow.TotalRowNo)
				SaveDocument(csvRow.TotalRowNo, &csvRow.Columns[1])
			}
			wg.Done()
		}()
	}
	wg.Wait()
	iig.Flush()
	log.Println("倒排索引建立完成")

	saveIDF()
	log.Println("IDF 建立完成")

	rows = ReadWukong()
	fenci.ReadIDF()
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func() {
			for csvRow := range rows {
				doc := strings.ToLower(csvRow.Columns[1])
				words, times := splitUniqueWords(&doc)
				// 存储正排索引
				topK := fenci.WeightTopK(doc, 10)
				topKWords := make([]string, len(topK))
				topKWeights := make([]float64, len(topK))
				for j := 0; j < len(topK); j++ {
					topKWords[j] = topK[j].Text
					topKWeights[j] = topK[j].Weight
				}
				AddWordsToForwardIndex(csvRow.TotalRowNo, words, times, topKWords, topKWeights)
				sig.AddWordsIdToSimpleInverted(topKWords, csvRow.TotalRowNo)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	sig.Flush()
	log.Println("正排索引建立完成")
	log.Println(docNum)
}

func InitTrie() {
	TrieTree = trie.NewTrie()
	for csvRow := range ReadWukong(1) {
		//构造trie树
		csvRow.Columns[1] = utils.RemoveSpace(csvRow.Columns[1])
		if len(csvRow.Columns) == 0 {
			continue
		}
		TrieTree.Add(csvRow.Columns[1])
	}
}

// splitUniqueWords 对文档进行分词，且记录关键词出现的次数
func splitUniqueWords(doc *string) ([]string, []int) {
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
	nodes := tree.InorderNode()
	words := make([]string, len(nodes))
	times := make([]int, len(nodes))
	for i := 0; i < len(nodes); i++ {
		words[i] = nodes[i].Val
		times[i] = nodes[i].Times
	}
	return words, times
}

func ReadWukong(num ...int) <-chan utils.CsvRow {
	if len(num) == 0 {
		files := make([]string, 10)
		for i := 0; i < 10; i++ {
			files[i] = utils.GetPath(fmt.Sprintf("/dataset/wukong_100m_%d.csv", i))
		}
		//return utils.ReadCsv(2, true, utils.GetPath("/dataset/wukong.csv"))
		return utils.ReadCsv(2, true, files...)
	} else {
		files := make([]string, num[0])
		for i := 0; i < num[0]; i++ {
			files[i] = utils.GetPath(fmt.Sprintf("/dataset/wukong_100m_%d.csv", i))
		}
		//return utils.ReadCsv(2, true, utils.GetPath("/dataset/wukong.csv"))
		return utils.ReadCsv(2, true, files...)
	}
}

// GetDocument 根据 ID 获取文档
func GetDocument(id int) ([]byte, bool) {
	key := []byte(strconv.Itoa(id))
	return storage.DocDB.Get(key)
}

// SaveDocument 保存文档至 leveldb 数据库
func SaveDocument(id int, doc *string) {
	err := storage.DocDB.Set([]byte(strconv.Itoa(id)), []byte(*doc))
	if err != nil {
		log.Println(doc, "存储失败")
		return
	}
}

// SaveWordId 将关键词和对应的文档 ID 存入 leveldb
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
		log.Fatalln(keyword, " SET 失败", err)
	}

	fmt.Printf("%d\n", docNum)
}

func saveIDF() {
	fileName := utils.GetPath("/database/idf.txt")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer file.Close()
	if err != nil {
		log.Println("idf file open failed", err)
		return
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	ch := storage.InvertedIndex.GetAllKeyValue()
	line := 0
	for kv := range ch {
		key := kv.Key()
		docs := storage.TransValueToIds(kv.Value())
		//docs := storage.InvertedIndex.GetDocStrIds([]byte(key))
		if len(docs) == 0 {
			return
		}

		// 计算逆文档频率 末尾加2是因为gse默认最小词频为2 且在使用idf的时候无法修改seg的配置
		// 在文档数据集比较小的时候经常出现 idf 小于 2 的情况，故只能出此下策
		idf := math.Log10(float64(docNum)/float64(len(docs))+1) + 2
		l := fmt.Sprintf("%s %f\n", key, idf)
		if l == "\n" {
			continue
		}
		_, err = writer.WriteString(l)
		if err != nil {
			log.Printf("写入idf失败: %s\n", err)
			return
		}
		line++
	}
	fenci.ReadIDF()
}

func shouldStoreIDF(word string) bool {
	runes := []rune(word)
	should := true
	for _, arune := range runes {
		if unicode.IsPunct(arune) || unicode.IsSpace(arune) || unicode.IsNumber(arune) || arune <= 64 {
			should = false
			break
		}
	}
	return should
}
