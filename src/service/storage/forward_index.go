package storage

import (
	"encoding/json"
	"log"
	"strconv"
)

type ForwardStore struct {
	Keywords    []string  `json:"keywords"`
	Times       []int     `json:"times"`
	TopKWords   []string  `json:"top_k_words"`
	TopKWeights []float64 `json:"top_k_weights"`
}

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

func (f ForwardIndexDBList) GetValueStruct(key []byte) ForwardStore {
	v, _ := f.Get(key)
	s := ForwardStore{}
	_ = json.Unmarshal(v, &s)
	return s
}

// Save 将关键词列表和次数列表存入数据库
func (store ForwardStore) Save(id int) {
	jsonVal, _ := json.Marshal(store)
	err := ForwardIndex.Set([]byte(strconv.Itoa(id)), jsonVal)
	if err != nil {
		log.Fatalln(id, jsonVal, " SET 失败", err)
	}
}
