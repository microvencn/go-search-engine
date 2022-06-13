package storage

import (
	"strconv"
	"strings"
)

type InvertedIndexDBList struct {
	shardDB
}

var InvertedIndex InvertedIndexDBList

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getInvertedIndexShardName(i))
	}
	InvertedIndex = InvertedIndexDBList{
		shardDB{
			Shard:  shard,
			DBList: dbs,
		},
	}
}

func (d InvertedIndexDBList) getShard(key []byte) int {
	if len(key) > 1 {
		return int(key[0]+key[1]) % d.Shard
	} else {
		return int(key[0]) % d.Shard
	}
}

func (d InvertedIndexDBList) GetDocStrIds(key []byte) []string {
	docs, exists := d.Get(key)
	if !exists {
		return make([]string, 0)
	}
	ids := strings.Split(string(docs), ",")
	return ids
}

func (d InvertedIndexDBList) GetDocIds(key []byte) []int {
	docs, exists := d.Get(key)
	if !exists {
		return make([]int, 0)
	}
	return TransValueToIds(string(docs))
}

// 获取分片数据库所在的路径
func getInvertedIndexShardName(shard int) string {
	return GetDBPath() + "/inverted_index_" + strconv.Itoa(shard)
}

func TransValueToIds(value string) []int {
	idStrings := strings.Split(value, ",")
	ids := make([]int, 0, len(idStrings)-1)
	for i := 0; i < len(idStrings)-1; i++ {
		id, err := strconv.Atoi(idStrings[i])
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
