package related_search

import (
	"go-search-engine/src/service/storage"
	"go-search-engine/src/service/utils"
	"log"
	"math"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"unicode"
)

var MinRelated float64 = 0.04
var MinDocNum int = 30

//var MaxDocNum int = 60

// 遍历多少词后将缓存写入数据库
var flushNum uint32 = 50000
var locks []*sync.Mutex
var num uint32 = 0

func InitStorage() {
	ch := storage.SimpleInvertedIndex.GetAllKeyValue()
	validCache := sync.Map{} // string sync.Map{string float64}
	//invalidCache := sync.Map{}
	//forwardCache := sync.Map{}  // int []string
	//invertedCache := sync.Map{} // string []int
	numCh := make(chan int)
	defer close(numCh)
	tooFewCache := sync.Map{}
	wg := sync.WaitGroup{}
	locks = make([]*sync.Mutex, 0, 5)

	go func() {
		for _ = range numCh {
			num++
			if num%1000 == 0 {
				log.Println("Progress:", num)
			}
			if num > flushNum {
				flushValidCache(&validCache)
			}
		}
	}()

	for i := 0; i < 6; i++ {
		locks = append(locks, &sync.Mutex{})
		wg.Add(1)
		i := i
		go func() {

			for kv := range ch {
				numCh <- 1
				locks[i].Lock()

				keyword := kv.Key()
				// 先判断是否应该保存，再判断是否存在于已经找到的文档数太少的缓存里
				if !shouldStore(keyword) {
					locks[i].Unlock()
					continue
				}
				if _, exists := tooFewCache.Load(keyword); exists {
					locks[i].Unlock()
					continue
				}

				// 获取关键词对应的文档
				// 若文档数量小于最小文档数量的限制 则加入缓存中保存
				keywordIds := storage.TransValueToIds(kv.Value())
				if len(keywordIds) < MinDocNum {
					tooFewCache.Store(keyword, true)
					locks[i].Unlock()
					continue
				}
				// 文档满足最少文档限制，则缓存查询结果
				//invertedCache.Store(keyword, keywordIds)

				// 取 MaxDocNum 限制范围内的文档 ID

				//keywordIds = keywordIds[0:utils.MinInt(MaxDocNum, len(keywordIds))]
				if len(keywordIds) < MinDocNum {
					locks[i].Unlock()
					continue
				}
				sort.Ints(keywordIds)
				words := make([]string, 0, 100)
				repeatedMap := make(map[string]bool)

				// 读取存在 keyword 的文章中的其它词
				// 并且无重复地加入 words 中
				for _, id := range keywordIds {
					aDocWords := storage.ForwardIndex.GetValueStruct([]byte(strconv.Itoa(id))).TopKWords
					//cacheDocs, _ := forwardCache.LoadOrStore(id, storage.ForwardIndex.GetValueStruct([]byte(strconv.Itoa(id))).TopKWords)
					//aDocWords := cacheDocs.([]string)
					for _, w := range aDocWords {
						if _, ok := repeatedMap[w]; ok {
							continue
						}
						words = append(words, w)
						repeatedMap[w] = true
					}
				}

				// 对每个词的文档转化为向量
				// 计算与 keyword 的相似度
				for _, word := range words {
					_, e := tooFewCache.Load(word)
					if e || !shouldStore(word) || word == keyword {
						continue
					}
					//var l, g string
					//if word < keyword {
					//	l = word
					//	g = keyword
					//} else {
					//	l = keyword
					//	g = word
					//}
					//if existsInCache(&invalidCache, l, g) {
					//	continue
					//}

					v, _ := storage.InvertedIndex.Get([]byte(word))
					wordIds := storage.TransValueToIds(string(v))
					sort.Ints(wordIds)
					// 判断文档数是否小于最小文档数
					if len(wordIds) < MinDocNum {
						tooFewCache.Store(word, true)
						continue
					}
					intersect := utils.CountIntersection(wordIds, keywordIds)
					var s float64
					if len(keywordIds)+len(wordIds)-intersect == 0 {
						s = 1
					} else {
						s = float64(intersect) / float64(len(keywordIds)+len(wordIds)-intersect)
					}
					////invertedCache.Store(word, keyword)
					//
					////docs := storage.SimpleInvertedIndex.GetDocIds([]byte(word))
					//wordIds = wordIds[0:utils.MinInt(len(wordIds), MaxDocNum)]
					//
					//s := cosSimilarity(transIdsToVector(keywordIds, wordIds))
					if s > MinRelated {
						storeInCache(&validCache, keyword, word, s)
						storeInCache(&validCache, word, keyword, s)
					}
					//else {
					//	storeInCache(&invalidCache, l, g, 0)
					//}
				}
				locks[i].Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	flushValidCache(&validCache)
}

func transIdsToVector(a []int, b []int) ([]int, []int) {
	// m key表示文章id value表示该id在向量中所在的维度
	m := make(map[int]int, len(a)+len(b))
	curIndex := 0
	for i := 0; i < len(a); i++ {
		if _, ok := m[a[i]]; !ok {
			m[a[i]] = curIndex
			curIndex += 1
		}
	}
	for i := 0; i < len(b); i++ {
		if _, ok := m[b[i]]; !ok {
			m[b[i]] = curIndex
			curIndex += 1
		}
	}
	v1 := make([]int, curIndex)
	v2 := make([]int, curIndex)
	for i := 0; i < len(a); i++ {
		v1[m[a[i]]] = 1
	}
	for i := 0; i < len(b); i++ {
		v2[m[b[i]]] = 1
	}
	return v1, v2
}

func cosSimilarity(v1 []int, v2 []int) float64 {
	var v1v2 int = 0
	for i := 0; i < len(v1); i++ {
		v1v2 += v1[i] * v2[i]
	}
	var absv1, absv2 int = 0, 0
	for i := 0; i < len(v1); i++ {
		absv1 += v1[i] * v1[i]
		absv2 += v2[i] * v2[i]
	}
	absv1v2 := math.Sqrt(float64(absv1)) * math.Sqrt(float64(absv2))
	return float64(v1v2) / absv1v2
}

func flushValidCache(m *sync.Map) {
	wg := sync.WaitGroup{}
	for i := 0; i < len(locks); i++ {
		i := i
		wg.Add(1)
		go func() {
			locks[i].Lock()
			wg.Done()
		}()
	}
	wg.Wait()
	m.Range(func(k1, v1 any) bool {
		wordList := make(storage.RelatedWords, 0, 10)
		w1 := k1.(string)
		m2 := v1.(*sync.Map)
		// 遍历获得 w1 w2 及其相似度 s
		m2.Range(func(k2, v2 any) bool {
			w2 := k2.(string)
			s := v2.(float64)
			rs := storage.RelatedWord{
				Text: w2,
				Sim:  s,
			}
			wordList = append(wordList, rs)
			return true
		})
		origin, _ := storage.Related.Get([]byte(w1))
		if origin != nil {
			wordList = append(origin, wordList...)
		}
		wordList = utils.RemoveRepeated(wordList)
		sort.Sort(wordList)
		// 倒序 成为降序
		for i := 0; i < len(wordList)/2; i++ {
			wordList[i], wordList[len(wordList)-1-i] = wordList[len(wordList)-1-i], wordList[i]
		}
		err := storage.Related.Set([]byte(w1), wordList)
		if err != nil {
			log.Println("Related set error:", err)
			return true
		}
		return true
	})
	*m = sync.Map{}
	num = 0
	for i := 0; i < len(locks); i++ {
		go locks[i].Unlock()
	}
	go runtime.GC()
}

func existsInCache(m *sync.Map, w1 string, w2 string) bool {
	v, exist := m.Load(w1)
	if !exist {
		return false
	}
	m2 := v.(*sync.Map)
	_, exist = m2.Load(w2)
	return exist
}

func storeInCache(m *sync.Map, w1 string, w2 string, s float64) {
	v, _ := m.LoadOrStore(w1, &sync.Map{})
	m2 := v.(*sync.Map)
	m2.Store(w2, s)
}

func shouldStore(word string) bool {
	runes := []rune(word)
	should := true
	if len(runes) < 2 {
		return false
	}
	for _, arune := range runes {
		if unicode.IsPunct(arune) || unicode.IsSpace(arune) || unicode.IsNumber(arune) ||
			arune <= 64 || arune >= 91 && arune <= 96 || arune >= 123 && arune <= 126 {
			should = false
			break
		} else {
			should = true
			break
		}
	}
	return should
}
