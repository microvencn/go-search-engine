package score

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func TestIntersect(t *testing.T) {
	//for time := 0; time < 10000; time++ {
	//	source, targets := getUniqueStringSlice(), getUniqueStringSlice()
	//	counter := Counter{
	//		TargetWords: targets,
	//	}
	//	result := counter.intersectSize(source)
	//
	//	m := make(map[string]int)
	//	nn := make([]string, 0)
	//	for _, v := range source {
	//		m[v]++
	//	}
	//	for _, v := range targets {
	//		times, _ := m[v]
	//		if times == 1 {
	//			nn = append(nn, v)
	//		}
	//	}
	//	count := len(nn)
	//
	//	if count != result {
	//		t.Errorf("输出了 %d，但是正确结果是 %d", result, count)
	//	} else {
	//		fmt.Println(time, ".", count, ".", result)
	//	}
	//}
}

func getUniqueStringSlice() []string {
	m := make(map[string]int)
	targets := make([]string, 100)
	for i := 0; i < 100; i++ {
		targets[i] = strconv.Itoa(rand.Int() % 600)
		if m[targets[i]] != 0 {
			i--
		} else {
			m[targets[i]] = 1
		}
	}
	sort.Strings(targets)
	return targets
}
