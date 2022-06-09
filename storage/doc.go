package storage

import "strconv"

var DocDB DocDBList

type DocDBList struct {
	DBList []*LevelDB
	Shard  int
}

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getDocShardName(i))
	}
	DocDB = DocDBList{
		Shard:  shard,
		DBList: dbs,
	}
}

func (d DocDBList) Get(key []byte) ([]byte, bool) {
	id, err := strconv.Atoi(string(key))
	if err != nil {
		return nil, false
	}
	shard := d.getShard(id)
	return d.DBList[shard].Get(key)
}

func (d DocDBList) Set(key []byte, value []byte) error {
	id, err := strconv.Atoi(string(key))
	if err != nil {
		return err
	}
	shard := d.getShard(id)
	err = d.DBList[shard].Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (d DocDBList) getShard(id int) int {
	return id % d.Shard
}

// 获取分片数据库所在的路径
func getDocShardName(shard int) string {
	return GetDBPath() + "/doc_db_" + strconv.Itoa(shard)
}
