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
	for word := range ch {
		count++
		fmt.Println(word)
	}
	fmt.Println(count)
}
