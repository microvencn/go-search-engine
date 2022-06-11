package fenci

type WordWeight struct {
	Text   string
	Weight float64
}

func (w WordWeight) Value() string {
	return w.Text
}

func (wws WordWeights) Len() int {
	return len(wws)
}

func (wws WordWeights) Less(i, j int) bool {
	return wws[i].Text < wws[j].Text
}

func (wws WordWeights) Swap(i, j int) {
	wws[i], wws[j] = wws[j], wws[i]
}

type WordWeights []WordWeight

func WeightTopK(text string, k int) WordWeights {
	tags := te.ExtractTags(text, k)
	//repeat := make(map[string]bool, k)

	// 获取 TOPK
	ww := make(WordWeights, 0, k)

	for i := 0; i < len(tags); i++ {
		weight := tags[i].Weight()
		//if weight == 0 {
		//	weight = 1
		//}
		ww = append(ww, WordWeight{
			Text:   tags[i].Text(),
			Weight: weight,
		})
		//repeat[tags[i].Text()] = true
	}
	return ww
}

func (wws *WordWeights) Normalize() {
	max := (*wws)[0].Weight
	min := (*wws)[len(*wws)-1].Weight
	for i, v := range *wws {
		(*wws)[i].Weight = wws.normalize(min, max, v.Weight)
	}
}

// 先进行最大最小归一化，然后结果 + 1
func (wws WordWeights) normalize(min float64, max float64, val float64) float64 {
	return (val-min)/(max-min) + 1
}

func Idf(key string) (float64, string, bool) {
	return te.Idf.Freq(key)
}

func (wws WordWeights) OnlyText() []string {
	texts := make([]string, len(wws))
	for i := 0; i < len(wws); i++ {
		texts[i] = wws[i].Text
	}
	return texts
}
