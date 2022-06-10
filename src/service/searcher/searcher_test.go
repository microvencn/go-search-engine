package searcher

import (
	"fmt"
	"go-search-engine/src/service/fenci"
	"testing"
)

func TestSimple(t *testing.T) {

	fenci.ReadDict()
	//idf, pos, exist := fenci.Idf("微信")
	//log.Println(idf, pos, exist)

	//s := make(fenci.WordWeights, 0, 10)
	//w := make([]string, 0, 10)
	//we := make([]float64, 0, 10)
	//for i := 0; i < 10; i++ {
	//	s = append(s, fenci.WordWeight{
	//		Text:   strconv.Itoa(i),
	//		Weight: 1,
	//	})
	//	w = append(w, strconv.Itoa(i))
	//	we = append(we, float64(1))
	//}
	//counter := score.Counter{
	//	TargetWords: s,
	//}
	//log.Println(counter.CosSimilarity(w, we))

	r := fenci.WeightTopK("百度图片", 10)
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i].Text, " ", r[i].Weight)
	}

	results := Simple("百度图片", 0, 20)
	for i := 0; i < len(results); i++ {
		fmt.Printf("%f %s\n", results[i].Score, results[i].Doc)
	}
	return
}
