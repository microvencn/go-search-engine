package keywords

import (
	"fmt"
	"testing"
)

func TestGetAllKeyWords(t *testing.T) {
	//fmt.Println(os.Getwd())
	//InitWukongIndex()
	ch := GetAllKeyWords()
	count := 0
	m := make(map[int]int)
	for word := range ch {
		bytes := []byte(word)
		total := 0
		for b := range bytes {
			total += b
		}
		m[total%10]++
		count++
		//fmt.Println(word)
	}
	for i, v := range m {
		fmt.Println(i, " ", v)
	}
	fmt.Println(count)
}
