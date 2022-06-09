package fenci

import (
	"fmt"
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/idf"
	"log"
)

var jieba gse.Segmenter

var te idf.TagExtracter

func ReadDict() {
	err := jieba.LoadDict()
	if err != nil {
		log.Fatalln("读取词典出错：", err)
		return
	}

	err = jieba.LoadStop()
	if err != nil {
		log.Fatalln("读取停止词词典出错：", err)
		return
	}

	te.WithGse(jieba)
	err = te.LoadIdf()
	if err != nil {
		fmt.Println("load idf: ", err)
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

func WeightTopK(text string, k int) idf.Segments {
	return te.ExtractTags(text, k)
}
