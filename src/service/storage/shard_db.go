package storage

import "strconv"

type ShardDBI interface {
	Get()
	Set()
	GetAllValue()
	GetAllKey()
}

// 分片 levelDB 的抽象结构体
// 实现了 Get 和 Set 方法，以及获取所有键或值的方法
// 在实践中一般只需要组合该结构体，然后重写 GetShard 方法即可
type shardDB struct {
	DBList []*LevelDB
	Shard  int
}

func (d shardDB) GetAllValue() <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < d.Shard; i++ {
			db := d.DBList[i].db
			iter := db.NewIterator(nil, nil)
			for iter.Next() {
				ch <- string(iter.Value())
			}
			iter.Release()
		}
		close(ch)
	}()
	return ch
}

func (d shardDB) GetAllKey() <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < d.Shard; i++ {
			db := d.DBList[i].db
			iter := db.NewIterator(nil, nil)
			for iter.Next() {
				ch <- string(iter.Value())
			}
			iter.Release()
		}
		close(ch)
	}()
	return ch
}

func (d shardDB) Get(key []byte) ([]byte, bool) {
	shard := d.GetShard(key)
	return d.DBList[shard].Get(key)
}

func (d shardDB) Set(key []byte, value []byte) error {
	shard := d.GetShard(key)
	return d.DBList[shard].Set(key, value)
}

// GetShard 获取指定键名所在的数据库切片
func (d shardDB) GetShard(key []byte) int {
	id, err := strconv.Atoi(string(key))
	if err != nil {
		return 0
	}
	return id % 10
}
