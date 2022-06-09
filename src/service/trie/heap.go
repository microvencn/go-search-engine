package trie

import (
	"container/heap"
	"fmt"
)

// RelatedWord 表示单词列表中的一个单词
type RelatedWord struct {
	Text  string // 单词
	Used  int    // 使用次数
	index int    // 索引
}

// WordList 定义单词列表，通常使用数组实现堆的存储
type WordList []*RelatedWord

//
// 实现以下5个函数以满足 heap.Interface 接口
//

func (wl WordList) Len() int {
	return len(wl)
}

func (wl WordList) Less(i, j int) bool {
	// 根据单词使用计数判断元素大小：单词使用计数越小，则元素值越小
	return wl[i].Used < wl[j].Used
}

func (wl WordList) Swap(i, j int) {
	wl[i], wl[j] = wl[j], wl[i]     // 交换元素
	wl[i].index, wl[j].index = i, j // 交换索引
}

func (wl *WordList) Push(x interface{}) {
	word := x.(*RelatedWord)
	word.index = len(*wl) // 根据接口说明，新元素的索引为Len()
	*wl = append(*wl, word)
}

func (wl *WordList) Pop() interface{} {
	// 根据接口说明，Pop()应移除并返回索引为Len()-1的元素
	word := (*wl)[len(*wl)-1]
	word.index = -1
	*wl = (*wl)[0 : len(*wl)-1]
	return word
}

// 该函数用来统计单词使用数量
func (wl *WordList) used(word string) {
	// 如果单词已在列表中，则增加使用计数
	i := 0
	for i = 0; i < len(*wl); i++ {
		if (*wl)[i].Text == word {
			(*wl)[i].Used++
			heap.Fix(wl, (*wl)[i].index) // 由于改变了索引为index的元素的大小，
			// 因此调用Fix重建堆
			break
		}
	}
	// 否则，添加单词到单词列表
	if i == len(*wl) {
		heap.Push(wl, &RelatedWord{
			Text:  word,
			Used:  1, // 使用计数为1
			index: -1,
		})
	}
}

// 限制单词列表长度，如果当前长度超过n，则将使用最少的单词移除，直到长度小于等于n为止
func (wl *WordList) limitTo(n int) {
	for wl.Len() > n {
		word := heap.Pop(wl).(*RelatedWord)
		fmt.Printf("Removed \"%s\" which has been used %d times\n",
			word.Text, word.Used)
	}
}

// 调试用，打印表示堆的切片
func (wl *WordList) debugPrint() {
	for i := 0; i < wl.Len(); i++ {
		fmt.Printf("\"%s\" at index %d used %d times\n",
			(*wl)[i].Text, (*wl)[i].index, (*wl)[i].Used)
	}
}

func main() {
	words := []string{
		"Gopher", "Tomcat", "Gopher", "Cynhard",
		"Gopher", "Tomcat", "Tomcat", "Gopher",
		"Go", "Go", "Heap"}

	wl := make(WordList, 0)

	// 使用任何堆操作之前，必须先进行初始化
	// 除非wl为空，或者已经排好序了
	heap.Init(&wl)

	for _, word := range words {
		wl.used(word)
	}

	wl.debugPrint()

	wl.limitTo(3)

	wl.debugPrint()

	wl.limitTo(0)
}
