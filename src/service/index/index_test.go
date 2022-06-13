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
	"time"
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
	r := storage.InvertedIndex.GetDocIds([]byte("a51"))
	for _, d := range r {
		doc, _ := GetDocument(d)
		fmt.Println(string(doc))
	}
	fmt.Println(len(r))
}

func TestTrie_Create(t *testing.T) {
	//前缀树构造逻辑：
	//在构造索引文件的时候同时构造前缀树。
	//构造前缀树时对分词结果进行自左向右的组合，例如‘字节跳动青训营’的分词结果为[字节 跳动 青训营],
	//则前缀树输入单词应该为[字节 跳动 青训营 字节跳动 字节跳动青训营 跳动青训营]
	trie := trie.NewTrie()
	trie.Add("字得其乐09 新世界")
	trie.Add("字母海鸥纹身")
	trie.Add("字母哥本身就是带着伤打球,而这一次他单脚落地真的太惊险了")
	trie.Add("字花萤光鞋帆布溜冰鞋成年双排滑轮旱冰鞋四轮")
	trie.Add("字体帮 境界")
	trie.Add("字母哥逆天空接引热议 詹皇抢镜坐场边表情落寞")
	trie.Add("字绣刺绣a绣针中格针三股绣11ct钝头50根针顶针拆线器")
	trie.Add("字针三角酪羊眼钉9钥匙扣超轻粘土橡皮泥彩泥软陶配件钥匙环加九")
	trie.Add("字符串或串是由数字,字母,下划线组成的一串字符.")
	trie.Add("字带一字扣罗马风露趾夏天凉鞋 t 韩版气质凉鞋女粗跟")
	trie.Add("字在其中字体/字形|可智奇 原创作品")
	trie.Add("字体帮 第600篇:荒芜 明日命题:风流")
	//trie.ToString()
	startTime := time.Now()
	wordList := trie.Search("字", -1)
	sort.Sort(wordList)
	elapsedTime := time.Since(startTime)
	fmt.Printf("搜索花费时间 %dms\n", elapsedTime)
	for _, v := range *wordList {
		fmt.Printf("%s %d\n", v.Text, v.Used)
	}
}

func TestTrie_Search(t *testing.T) {
	//构造前缀树
	fenci.ReadDict()
	InitTrie()
	//搜索提示词
	//输入关键词返回与关键词相关的字符串，返回出现次数最高的几个相关字符串
	startTime := time.Now()
	wordList := TrieTree.Search("字", -1)
	if wordList != nil {
		sort.Sort(wordList)
		elapsedTime := time.Since(startTime)
		fmt.Printf("搜索花费时间 %dms\n", elapsedTime)
		//fmt.Println(wordList)
		for _, v := range *wordList {
			fmt.Println(v)
		}
	} else {
		fmt.Println("结果为空")
	}
}
