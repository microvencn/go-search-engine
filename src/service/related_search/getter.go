package related_search

import (
	"go-search-engine/src/service/fenci"
	"go-search-engine/src/service/storage"
	"go-search-engine/src/service/utils"
	"sort"
)

func GetWordRelatedWords(word string) (storage.RelatedWords, bool) {
	return storage.Related.Get([]byte(word))
}

func GetQueryRelatedWords(query string) storage.RelatedWords {
	topk := fenci.WeightTopK(query, 3)
	words := make(storage.RelatedWords, 0, 20)
	for _, word := range topk.OnlyText() {
		related, e := GetWordRelatedWords(word)
		if e {
			words = append(words, related...)
		}
	}
	utils.RemoveRepeated(words)
	sort.Sort(words)
	for i := 0; i < len(words)/2; i++ {
		words[i], words[len(words)-1-i] = words[len(words)-1-i], words[i]
	}
	return words
}
