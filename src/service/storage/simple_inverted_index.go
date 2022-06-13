package storage

import "strconv"

type SimpleInvertedIndexDBList struct {
	InvertedIndexDBList
}

var SimpleInvertedIndex SimpleInvertedIndexDBList

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getSimpleInvertedIndexShardName(i))
	}
	SimpleInvertedIndex = SimpleInvertedIndexDBList{
		InvertedIndexDBList{
			shardDB{
				Shard:  shard,
				DBList: dbs,
			},
		},
	}
}

func getSimpleInvertedIndexShardName(shard int) string {
	return GetDBPath() + "/simple_inverted_index_" + strconv.Itoa(shard)
}
