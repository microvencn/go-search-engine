package storage

import "strconv"

type ForwardIndexDBList struct {
	DocDBList
}

var ForwardIndex ForwardIndexDBList

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getForwardIndexShardName(i))
	}
	ForwardIndex = ForwardIndexDBList{
		DocDBList{
			shardDB{
				Shard:  10,
				DBList: dbs,
			},
		},
	}
}

func getForwardIndexShardName(shard int) string {
	return GetDBPath() + "/forward_index_" + strconv.Itoa(shard)
}
