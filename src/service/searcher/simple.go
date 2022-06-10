package searcher

import (
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/index"
	"go-search-engine/src/service/score"
	"go-search-engine/src/service/storage"
	"go-search-engine/src/service/utils"
	"sort"
	"strings"
)

type simple struct {
	Doc string
	score.IdScore
}

type SimpleResult []simple

// Simple 简单搜索，offset 搜索结果的偏移，length 为需要的结果数量
func Simple(query string, offset int, length int) SimpleResult {
	query = strings.ToLower(query)
	// 分词并保存至 targets，将所有文档 id 存储至 ids
	targets := make(fenci.WordWeights, 0, len(query)/2)
	ids := make([]int, 0)
	repeatedMap := make(map[string]bool, len(query))

	// 对用户输入尝试取 TOPK 并去重
	words := fenci.WeightTopK(query, 10)
	for i := 0; i < len(words); i++ {
		if !repeatedMap[words[i].Text] {
			targets = append(targets, words[i])
			repeatedMap[words[i].Text] = true
			id, _ := index.GetWordIds(words[i].Text)
			ids = append(ids, id...)
		}
	}
	sort.Sort(targets)

	if len(words) < 10 {
		fenci.ExecAndDoSomething(&query, func(word string) {
			if !repeatedMap[word] {
				targets = append(targets, fenci.WordWeight{
					Text:   word,
					Weight: 1,
				})
				repeatedMap[word] = true
				id, _ := index.GetWordIds(word)
				ids = append(ids, id...)
			}
		})
	}
	ids = utils.RemoveRepeatedElement(ids)

	// 初始化分数计算器，将用户输入的分词结果作为分数计算依据
	// 对 ids 中对应的所有文档进行排序
	counter := score.Counter{
		TargetWords: targets,
	}
	idScores := counter.SortAfterCount(ids)

	// 生成搜索结果
	results := make([]simple, utils.MinInt(length, len(idScores)))
	size := length
	for i := len(idScores) - 1 - offset; size > 0 && i > -1; i-- {
		idScore := idScores[i]
		doc, _ := storage.GetDocument(idScore.Id)
		results[length-size] = simple{
			Doc:     string(doc),
			IdScore: idScore,
		}
		size--
	}
	return results
}
