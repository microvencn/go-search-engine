package fenci

import (
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/idf"
	"go-search-engine/src/service/utils"
	"log"
)

var jieba gse.Segmenter

var te idf.TagExtracter

func ReadDict() {
	err := jieba.LoadDict()
	jieba.MinTokenFreq = 0
	//jieba.Load
	//jieba.AddToken("字节跳动", 3.5, "n")
	if err != nil {
		log.Fatalln("读取词典出错：", err)
		return
	}
	//a, b, c := jieba.Dict.Find([]byte("图片"))
	//fmt.Println(a, b, c)

	ReadIDF()

	err = jieba.LoadStop()
	if err != nil {
		log.Fatalln("读取停止词词典出错：", err)
		return
	}
}

func doSomething(words []string, f func(word string)) {
	for _, word := range words {
		f(word)
	}
}

// ExecAndDoSomething 对 sentence 分词并对每个词执行 f 函数
func ExecAndDoSomething(sentence *string, f func(word string)) {
	doSomething(jieba.CutSearch(*sentence, true), f)
}

func Cut(doc string) []string {
	return jieba.CutSearch(doc, true)
}

func ReadIDF() {
	te.WithGse(jieba)
	err := te.LoadIdf(utils.GetPath("/database/idf.txt"))

	if err != nil {
		//log.Fatalln("加载 IDF 词典失败: ", err)
		return
	}
	_ = te.Idf.AddToken("go-search-engine", 1, "")
	return
}
