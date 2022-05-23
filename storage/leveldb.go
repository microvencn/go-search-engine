package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db   *leveldb.DB
	path string
}

var DocDB *LevelDB
var DictDB *LevelDB

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
	DictDB = Open("./database/dict")
	DocDB = Open("./database/doc")
}
