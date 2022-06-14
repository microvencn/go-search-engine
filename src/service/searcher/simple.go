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
	Doc string `json:"doc"`
	*score.IdScore
}

type SimpleResult []simple

func simpleSearch(query string) score.IdScoreList {
	query = strings.ToLower(query)
	// 分词并保存至 targets，将所有文档 id 存储至 ids
	targets := make(fenci.WordWeights, 0, len(query)/2)

	repeatedMap := make(map[string]bool, len(query))

	// 对用户输入尝试取 TOPK 并去重
	queryTokens := fenci.WeightTopK(query, utils.MaxInt(len(query)/2, 1))
	ids := make([]int, 0, len(queryTokens)*5)
	for i := 0; i < len(queryTokens); i++ {
		if _, ok := repeatedMap[queryTokens[i].Text]; !ok {
			targets = append(targets, queryTokens[i])
			repeatedMap[queryTokens[i].Text] = true
			id, _ := index.GetSimpleWordIds(queryTokens[i].Text)
			ids = append(ids, id...)
		}
	}
	sort.Sort(targets)

	// 关键词数量为0时，取停止词
	if len(ids) == 0 {
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
		ids = utils.RemoveRepeatedElement(ids)
		idScores := make(score.IdScoreList, 0, len(ids))
		for i := 0; i < len(ids); i++ {
			idScores = append(idScores, &score.IdScore{
				Id:    ids[i],
				Score: 0,
			})
		}
		return idScores
	}
	ids = utils.RemoveRepeatedElement(ids)

	// 初始化分数计算器，将用户输入的分词结果作为分数计算依据
	// 对 ids 中对应的所有文档进行排序
	counter := score.Counter{
		TargetWords: targets,
	}
	idScores := counter.SortAfterCount(ids)

	return idScores
}

func idScoresToSimpleResult(idScores score.IdScoreList) SimpleResult {
	results := make(SimpleResult, 0, len(idScores))
	for i := 0; i < len(idScores); i++ {
		idScore := idScores[i]
		doc, exists := storage.GetDocument(idScore.Id)
		if !exists {
			continue
		}
		results = append(results, simple{
			Doc:     string(doc),
			IdScore: idScore,
		})
	}
	return results
}

// Simple 简单搜索，offset 搜索结果的偏移，length 为需要的结果数量
func Simple(query string, offset int, length int) (int, SimpleResult) {
	if offset < 0 {
		return 0, nil
	}
	idScores := simpleSearch(query)
	// 生成搜索结果
	return len(idScores), idScoresToSimpleResult(idScores[offset:utils.MinInt(length, len(idScores)-offset)])
}

func SimpleWithFilter(query string, offset int, length int, filter []string) (int, SimpleResult) {
	if offset < 0 {
		return 0, nil
	}
	idScores := simpleSearch(query)

	// 生成搜索结果
	results := make(score.IdScoreList, 0, utils.MaxInt(length, 0))
	size := length
	// 过滤搜索结果
	if filter != nil {
		sort.Strings(filter)
	}
	for i := offset; size > 0 && i < len(idScores); i++ {
		idScore := idScores[i]
		if filter != nil {
			words, exists := index.GetIdWords(idScore.Id)
			if !exists {
				continue
			}
			if utils.HasIntersection(words.Keywords, filter) {
				continue
			}
		}
		results = append(results, idScores[i])
		size--
	}
	return len(idScores), idScoresToSimpleResult(results)
}
