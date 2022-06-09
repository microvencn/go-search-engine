package fenci

import (
	"fmt"
	"testing"
)

func TestFenci(t *testing.T) {
	ReadDict()
	doc := "百度啥百度大佬你到底会不会啊"
	ExecAndDoSomething(&doc, func(word string) {
		fmt.Println(word)
	})
}
