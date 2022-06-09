package storage

import "strconv"

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
		return int(key[0]+key[1]) % 10
	} else {
		return int(key[0]) % 10
	}
}

// 获取分片数据库所在的路径
func getInvertedIndexShardName(shard int) string {
	return GetDBPath() + "/inverted_index_" + strconv.Itoa(shard)
}
