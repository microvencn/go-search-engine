package score

import (
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/index"
	"math"
	"sort"
	"sync"
)

type IdScore struct {
	Id    int     `json:"id"`
	Score float64 `json:"score"`
}
type IdScoreList []*IdScore

type Counter struct {
	TargetWords fenci.WordWeights
}

func (c Counter) CountById(id int) (score float64, success bool) {
	forward, exists := index.GetIdWords(id)
	if !exists {
		return 0, false
	}
	return c.CosSimilarity(forward.TopKWords, forward.TopKWeights), true
}

func (c Counter) SortAfterCount(ids []int) IdScoreList {
	idScoreCh := make(chan *IdScore)
	idCh := make(chan int)

	go func() {
		for i := 0; i < len(ids); i++ {
			idCh <- ids[i]
		}
		close(idCh)
	}()

	countWg := sync.WaitGroup{}
	for i := 0; i < 6; i++ {
		countWg.Add(1)
		go func() {
			for id := range idCh {
				score, success := c.CountById(id)
				if !success {
					continue
				}
				idScoreCh <- &IdScore{
					Score: score,
					Id:    id,
				}
			}
			countWg.Done()
		}()
	}

	go func() {
		countWg.Wait()
		close(idScoreCh)
	}()

	idScores := make(IdScoreList, 0, len(ids)*2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for idScore := range idScoreCh {
			idScores = append(idScores, idScore)
		}
		wg.Done()
	}()
	wg.Wait()

	sort.Sort(idScores)
	reverse(idScores)
	return idScores
}

// Intersection 用交集大小计算指定关键词的得分，要求 targetWords 和 words 均为有序
func (c Counter) Intersection(words []string, times []int) float64 {
	i, j := 0, 0
	var count float64 = 0
	for {
		for i < len(c.TargetWords) && c.TargetWords[i].Text < words[j] {
			i++
		}
		if i == len(c.TargetWords) {
			break
		}

		for j < len(words) && words[j] < c.TargetWords[i].Text {
			j++
		}
		if j == len(words) {
			break
		}

		if c.TargetWords[i].Text == words[j] {
			count++
			i++
			j++
			if i == len(c.TargetWords) || j == len(words) {
				break
			}
		}
	}
	return count
}

// CosSimilarity 计算TopK 的余弦相似度
func (c Counter) CosSimilarity(words []string, weights []float64) float64 {
	length := len(words) + len(c.TargetWords)
	// m 存储单词在向量中对应的维度
	m := make(map[string]int, length)
	vIndex := 0
	for _, word := range words {
		m[word] = vIndex
		vIndex += 1
	}
	for _, word := range c.TargetWords {
		if _, ok := m[word.Text]; ok {
			continue
		}
		m[word.Text] = vIndex
		vIndex += 1
	}

	v1, v2 := make([]float64, length), make([]float64, length)
	for i := 0; i < len(words); i++ {
		word := words[i]
		v1[m[word]] = weights[i]
	}
	for _, word := range c.TargetWords {
		v2[m[word.Text]] = word.Weight
	}

	var v1v2 float64 = 0
	for i := 0; i < length; i++ {
		v1v2 += v1[i] * v2[i]
	}
	var absv1, absv2 float64 = 0, 0
	for i := 0; i < length; i++ {
		absv1 += v1[i] * v1[i]
		absv2 += v2[i] * v2[i]
	}
	absv1 = math.Sqrt(absv1)
	absv2 = math.Sqrt(absv2)
	return v1v2 / (absv2 * absv1)
}

func (p IdScoreList) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}

func (p IdScoreList) Len() int {
	return len(p)
}

func (p IdScoreList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func reverse(s IdScoreList) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
