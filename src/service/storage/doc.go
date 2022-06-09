package storage

import (
	"strconv"
)

var DocDB DocDBList

type DocDBList struct {
	shardDB
}

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getDocShardName(i))
	}
	DocDB = DocDBList{
		shardDB{
			Shard:  shard,
			DBList: dbs,
		},
	}
}

// 获取分片数据库所在的路径
func getDocShardName(shard int) string {
	return GetDBPath() + "/doc_db_" + strconv.Itoa(shard)
}

func (d DocDBList) GetAllDoc() <-chan string {
	return d.GetAllValue()
}

// GetDocument 根据 ID 获取文档
func GetDocument(id int) ([]byte, bool) {
	key := []byte(strconv.Itoa(id))
	return DocDB.Get(key)
}
