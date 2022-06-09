package searcher

import (
	"fmt"
	"go-search-engine/src/service/fenci"
	"testing"
)

func TestSimple(t *testing.T) {
	fenci.ReadDict()
	r := fenci.WeightTopK("好", 10)
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i].Text(), " ", r[i].Weight())
	}
	//fenci.ReadDict()
	//results := Simple("好用", 0, 10)
	//for i := 0; i < len(results); i++ {
	//	log.Printf("%f %s", results[i].Score, results[i].Doc)
	//}
	//return
}
