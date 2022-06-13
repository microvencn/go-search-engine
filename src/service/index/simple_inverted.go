package index

import (
	"bytes"
	"go-search-engine/src/service/storage"
	"log"
	"strconv"
	"strings"
	"sync"
)

type SimpleInvertedGen struct {
	cache     sync.Map
	cacheSize uint32
}

var sig SimpleInvertedGen

func (s *SimpleInvertedGen) AddWordsIdToSimpleInverted(words []string, id int) {
	for _, word := range words {
		store, _ := (*s).cache.LoadOrStore(word, &bytes.Buffer{})
		buffer, _ := store.(*bytes.Buffer)
		(*buffer).WriteString(strconv.Itoa(id) + ",")
		(*s).cache.Store(word, buffer)
	}
}

func (s *SimpleInvertedGen) Flush() {
	sig.cache.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*bytes.Buffer)
		err := storage.SimpleInvertedIndex.Set([]byte(k), (*v).Bytes())
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	})
}

func GetSimpleWordIds(word string) ([]int, bool) {
	idByteList, e := storage.SimpleInvertedIndex.Get([]byte(word))
	if !e {
		return nil, false
	}
	idList := string(idByteList)

	// 由于存储时使用逗号作为分隔。所以读取时使用逗号分割
	idStr := strings.Split(idList, ",")
	// 这里减一是因为存入数据库时以 , 结尾，idStr 的末尾元素为空串
	ids := make([]int, 0, len(idStr)-1)
	for i := 0; i < len(idStr)-1; i++ {
		id, _ := strconv.Atoi(idStr[i])
		ids = append(ids, id)
	}
	return ids, true
}
