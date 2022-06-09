package trie

import (
	"container/heap"
	"fmt"
)

// 字典树：
/*
	1. 单词结尾表示：使用一个int类型，表示以当前节点为结尾字符的单词数量
		叶子结点是单词结尾
		非叶子节点也可能是单词的结尾
	2. 实现功能
		增加单词
		搜索单词
*/

// Node 定义树的节点
type Node struct {
	//  若干个子节点: 字符：节点
	Children map[rune]*Node
	//  表示以该节点为结尾字符的单词数量（用来做相关搜索）
	cntWord int
}

// Add 给子节点 列表添加新的节点
func (n *Node) Add(c rune) {
	// 构造节点：nil类型map不可直接使用添加元素等
	node := Node{
		Children: make(map[rune]*Node, 0),
	}
	n.Children[c] = &node
}

// Trie 整棵树的根节点
type Trie struct {
	// 字典树的根节点
	Root *Node
}

// InitTrie Init 初始化根节点
func (t *Trie) InitTrie() *Trie {
	t.Root = &Node{
		make(map[rune]*Node, 0),
		0,
	}
	return t
}

// NewTrie New 外部调用初始化字典树
func NewTrie() *Trie { return new(Trie).InitTrie() }

// Add 添加单词
func (t *Trie) Add(word string) {
	// 找到每一个节点
	cur := t.Root
	// 遍历word的每一个字符
	for _, c := range word {
		// 1.查找子节点中是否包含当前字符
		// 如果不包含,添加一下
		if _, ok := cur.Children[c]; !ok {
			// 添加新节点
			cur.Add(c)
		}
		// cur指向下一个子节点
		cur = cur.Children[c]
	}
	// 将最后一个节点字符的数量+1
	cur.cntWord++
}

// Contains 判断是否包含指定的单词
func (t *Trie) Contains(word string) *Node {
	// 找到每一个节点
	cur := t.Root
	// 遍历word的每一个字符
	for _, c := range word {
		// 1. 查找子节点中是否包含当前字符
		if _, ok := cur.Children[c]; !ok {
			return nil
		}
		// cur 指向下一个子节点
		cur = cur.Children[c]
	}
	// 判断最后一个单词是否是单词末尾 cur.cntWord > 0,
	return cur
}

type queueNode struct {
	node *Node
	path []rune
}

func (t *Trie) ToString() {
	queue := make([]*Node, 0)
	queue = append(queue, t.Root)
	for len(queue) > 0 {
		q := queue[0]
		queue = queue[1:]
		for key, node := range q.Children {
			fmt.Printf("%c ", key)
			queue = append(queue, node)
		}
		fmt.Println()
	}
}

// Search 查询以word为前缀的limit个字符串
func (t *Trie) Search(word string, limit int) *WordList {
	cur := t.Contains(word)
	if cur == nil {
		return nil
	}
	wl := make(WordList, 0)
	heap.Init(&wl)
	queue := make([]queueNode, 0)
	queue = append(queue, queueNode{cur, []rune(word)})
	for len(queue) > 0 && (limit == -1 || len(wl) < limit) {
		q := queue[0]
		queue = queue[1:]
		for key, node := range q.node.Children {
			path := make([]rune, len(q.path))
			copy(path, q.path)
			path = append(path, key)
			if node.cntWord > 0 {
				wl.Push(&RelatedWord{string(path), node.cntWord, -1})
			}
			queue = append(queue, queueNode{node, path})

		}
	}
	return &wl
}

//func dfs(relatedWord *[]RelatedWord, cur *Node, prefix *[]rune, cnt int, limit int) {
//	if cur == nil || cnt > limit {
//		return
//	} else if cur.cntWord > 0 {
//
//	}
//
//	for _, c := range cur.Children {
//
//	}
//}
