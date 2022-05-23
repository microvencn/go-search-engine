package fenci

import (
	"github.com/wangbin/jiebago"
	"log"
)

var seg jiebago.Segmenter

func ReadDict() {
	err := seg.LoadDictionary("./fenci/dictionary.txt")
	if err != nil {
		log.Fatalln("加载词典失败")
		return
	}
}

func doSomething(ch <-chan string, f func(word string)) {
	for word := range ch {
		f(word)
	}
}

// ExecAndDoSomething 对 sentence 分词并对每个词执行 f 函数
func ExecAndDoSomething(sentence string, f func(word string)) {
	doSomething(seg.CutForSearch(sentence, true), f)
}
