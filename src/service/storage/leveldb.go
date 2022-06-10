package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"go-search-engine/src/service/utils"
)

type LevelDB struct {
	db   *leveldb.DB
	path string
}

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

func GetDBPath() string {
	return utils.GetPath("/database")
}
