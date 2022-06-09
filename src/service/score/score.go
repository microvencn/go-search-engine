package score

import (
	"github.com/go-ego/gse/hmm/idf"
	"go-search-engine/src/service/index"
	"sort"
)

type IdScore struct {
	Id    int
	Score float64
}
type IdScoreList []IdScore

type Counter struct {
	TargetWords idf.Segments
}

func (c Counter) CountById(id int) (score float64, success bool) {
	words, exists := index.GetIdWords(id)
	if !exists {
		return 0, false
	}
	return c.wordsScore(words.Keywords, words.Times), true
}

func (c Counter) SortAfterCount(ids []int) IdScoreList {
	scores := make(IdScoreList, len(ids))
	for i := 0; i < len(ids); i++ {
		score, _ := c.CountById(ids[i])
		scores[i] = IdScore{
			Score: score,
			Id:    ids[i],
		}
	}
	sort.Sort(scores)
	return scores
	//sorted := make([]int, len(ids))
	//for i := 0; i < len(ids); i++ {
	//	sorted[i] = scores[i].id
	//}
	//return sorted
}

// wordsScore 计算指定关键词的得分，要求 targetWords 和 words 均为有序
func (c Counter) wordsScore(words []string, times []int) float64 {
	i, j := 0, 0
	var count float64 = 0
	for {
		for i < len(c.TargetWords) && c.TargetWords[i].Text() < words[j] {
			i++
		}
		if i == len(c.TargetWords) {
			break
		}

		for j < len(words) && words[j] < c.TargetWords[i].Text() {
			j++
		}
		if j == len(words) {
			break
		}

		if c.TargetWords[i].Text() == words[j] {
			count += float64(times[j]) + c.TargetWords[i].Weight()*float64(times[j])
			i++
			j++
			if i == len(c.TargetWords) || j == len(words) {
				break
			}
		}
	}
	return count
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
