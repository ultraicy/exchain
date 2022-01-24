package mpt

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const (
	EvmDataDir = "data"
	EvmSpace   = "evm_mpt"

	FlagDBBackend = "db_backend"

	FlagTrieDirtyDisabled = "trie-dirty-disabled"
	FlagTrieCacheSize     = "trie-cache-size"
)

var (
	TrieDirtyDisabled      = false
	TrieCacheSize     uint = 4096 // MB
)

func InstanceOfEvmStore() ethstate.Database {
	homeDir := viper.GetString(flags.FlagHome)
	path := filepath.Join(homeDir, EvmDataDir)

	backend := viper.GetString(FlagDBBackend)
	if backend == "" {
		backend = string(dbm.GoLevelDBBackend)
	}

	kvstore, e := CreateKvDB(EvmSpace, BackendType(backend), path)
	if e != nil {
		panic("fail to open database: " + e.Error())
	}
	db := rawdb.NewDatabase(kvstore)
	gEvmMptDatabase := ethstate.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     int(TrieCacheSize),
		Journal:   "",
		Preimages: true,
	})

	return gEvmMptDatabase
}

type BackendType string

// These are valid backend types.
const (
	// GoLevelDBBackend represents goleveldb (github.com/syndtr/goleveldb - most
	// popular implementation)
	//   - pure go
	//   - stable
	GoLevelDBBackend BackendType = "goleveldb"

	// RocksDBBackend represents rocksdb (uses github.com/tecbot/gorocksdb)
	//   - EXPERIMENTAL
	//   - requires gcc
	//   - use rocksdb build tag (go build -tags rocksdb)
	RocksDBBackend BackendType = "rocksdb"
)

func init() {
	dbCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewLevelDB(name, dir)
	}
	registerDBCreator(GoLevelDBBackend, dbCreator, false)
}

func NewLevelDB(name string, dir string) (ethdb.KeyValueStore, error) {
	file := filepath.Join(dir, name+".db")
	return leveldb.New(file, 128, 1024, EvmSpace, false)
}

type dbCreator func(name string, dir string) (ethdb.KeyValueStore, error)

var backends = map[BackendType]dbCreator{}

func registerDBCreator(backend BackendType, creator dbCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

func CreateKvDB(name string, backend BackendType, dir string) (ethdb.KeyValueStore, error) {
	dbCreator, ok := backends[backend]
	if !ok {
		keys := make([]string, len(backends))
		i := 0
		for k := range backends {
			keys[i] = string(k)
			i++
		}
		panic(fmt.Sprintf("Unknown db_backend %s, expected either %s", backend, strings.Join(keys, " or ")))
	}

	return dbCreator(name, dir)
}

var (
	KeyPrefixLatestHeight       = []byte{0x01}
	KeyPrefixRootMptHash        = []byte{0x02}
	KeyPrefixLatestStoredHeight = []byte{0x03}
)

// GetLatestBlockHeight get latest mpt storage height
func GetLatestBlockHeight(db ethstate.Database) uint64 {
	rst, err := db.TrieDB().DiskDB().Get(KeyPrefixLatestHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// GetRootMptHash gets root mpt hash from block height
func GetRootMptHash(db ethstate.Database, height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := db.TrieDB().DiskDB().Get(append(KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}
