package related_search

import (
	"go-search-engine/src/service/fenci"
	"log"
	"testing"
)

func TestInitRelated(t *testing.T) {
	InitStorage()
}

func TestGetRelatedWord(t *testing.T) {
	fenci.ReadDict()
	//ch := storage.Related.GetAllKeyValue()
	//for kv := range ch {
	//	fmt.Print(kv.Key())
	//	rs, _ := storage.Related.Get([]byte(kv.Key()))
	//	for _, w := range rs {
	//		fmt.Printf(" %s %.2f", w.Text, w.Sim)
	//	}
	//	fmt.Print("\n")
	//}
	//log.Println(GetWordRelatedWords("姜子牙"))
	log.Println(GetQueryRelatedWords("百度"))
}
