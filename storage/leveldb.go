package storage

import (
	"GoSearchEngine/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

type LevelDB struct {
	db   *leveldb.DB
	path string
}

var InvertedIndex *LevelDB
var ForwardIndex *LevelDB

func Open(path string) *LevelDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil
	}
	return &LevelDB{
		db:   db,
		path: path,
	}
}

func (s *LevelDB) Get(key []byte) ([]byte, bool) {

	buffer, err := s.db.Get(key, nil)
	if err != nil {
		return nil, false
	}
	return buffer, true
}

func (s *LevelDB) Set(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *LevelDB) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *LevelDB) Close() error {
	return s.db.Close()
}

func init() {
	InvertedIndex = Open(utils.GetPath("/database/inverted"))
	ForwardIndex = Open(utils.GetPath("/database/forward"))
}

// GetDocument 根据 ID 获取文档
func GetDocument(id int) ([]byte, bool) {
	key := []byte(strconv.Itoa(id))
	return DocDB.Get(key)
}

func GetDBPath() string {
	return utils.GetPath("/database")
}
