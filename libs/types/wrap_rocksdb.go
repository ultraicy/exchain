//go:build rocksdb
// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/pkg/errors"
	"github.com/tecbot/gorocksdb"
)

func init() {
	dbCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewWrapRocksDB(name, dir)
	}
	registerDBCreator(RocksDBBackend, dbCreator, false)
}

type WrapRocksDB struct {
	*tmdb.RocksDB
}

func NewWrapRocksDB(name string, dir string) (*WrapRocksDB, error) {
	rdb, err := tmdb.NewRocksDB(name, dir)

	return &WrapRocksDB{rdb}, err
}

func (db *WrapRocksDB) Put(key []byte, value []byte) error {
	return db.Set(key, value)
}

func (db *WrapRocksDB) NewBatch() ethdb.Batch {
	return NewWrapRocksDBBatch(db.RocksDB)
}

func (db *WrapRocksDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	st := append(prefix, start...)
	return NewWrapRocksDBIterator(db.RocksDB, st, nil)
}

func (db *WrapRocksDB) Stat(property string) (string, error) {
	stats := db.RocksDB.Stats()
	if pro, ok := stats[property]; ok {
		return pro, nil
	}

	return "", errors.New("property not exist")
}

func (db *WrapRocksDB) Compact(start []byte, limit []byte) error {
	db.DB().CompactRange(gorocksdb.Range{Start: start, Limit: limit})
	return nil
}
