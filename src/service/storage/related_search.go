package storage

import (
	"encoding/json"
	"strconv"
)

type RelatedDBList struct {
	shardDB
}

type RelatedWord struct {
	Text string  `json:"text"`
	Sim  float64 `json:"sim"`
}

type RelatedWords []RelatedWord

func (r RelatedWords) Less(i, j int) bool {
	return r[i].Sim < r[j].Sim
}

func (r RelatedWords) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RelatedWords) Len() int {
	return len(r)
}

var Related RelatedDBList

func init() {
	shard := 10
	dbs := make([]*LevelDB, shard)
	for i := 0; i < shard; i++ {
		dbs[i] = Open(getRelatedShardName(i))
	}
	Related = RelatedDBList{
		shardDB{
			Shard:  shard,
			DBList: dbs,
		},
	}
}

func (d RelatedDBList) getShard(key []byte) int {
	if len(key) > 1 {
		return int(key[0]+key[1]) % d.Shard
	} else {
		return int(key[0]) % d.Shard
	}
}

func (d RelatedDBList) Get(key []byte) (RelatedWords, bool) {
	v, exists := d.DBList[d.getShard(key)].Get(key)
	if !exists {
		return nil, false
	}
	s := make(RelatedWords, 0)
	_ = json.Unmarshal(v, &s)
	return s, true
}

func (d RelatedDBList) Set(key []byte, val RelatedWords) error {
	v, _ := json.Marshal(val)
	shard := d.getShard(key)
	return d.DBList[shard].Set(key, v)
}

func (w RelatedWord) Value() string {
	return w.Text
}

// 获取分片数据库所在的路径
func getRelatedShardName(shard int) string {
	return GetDBPath() + "/related_" + strconv.Itoa(shard)
}
