package avl_struct

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var num int = 0
var mu sync.Mutex

func TestAvl(t *testing.T) {
	wg := sync.WaitGroup{}
	// 10 * 10000 = 1e5 次
	for i := 0; i < 10; i++ {
		go test(&wg, t)
		wg.Add(1)
	}
	wg.Wait()
	return
}

// 二叉平衡数测试
func test(wg *sync.WaitGroup, t *testing.T) {
	avl := Init[int]()

	//测试 10000 次
	for i := 0; i < 10000; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// 生成 1000 个数
		for i := 0; i < 999; i++ {
			avl.Insert(r.Intn(1000))
		}
		// 校验先序遍历是否有序
		s := avl.PreOrder()
		for i := 1; i < len(s); i++ {
			if s[i] < s[i-1] {
				t.Errorf("存在错误")
				return
			}
		}
		mu.Lock()
		num++
		fmt.Printf("第 %d 次测试正确 \n", num)
		mu.Unlock()
	}
	wg.Done()
}
