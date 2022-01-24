package mpt

import (
	"encoding/binary"
	"github.com/okex/exchain/libs/types"
	"path/filepath"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const (
	AccDataDir = "data"
	AccSpace   = "acc"

	FlagDBBackend = "db_backend"
)

var (
	gAccMptDatabase ethstate.Database = nil
	initOnce        sync.Once
)

func InstanceOfAccStore() ethstate.Database {
	initOnce.Do(func() {
		homeDir := viper.GetString(flags.FlagHome)
		path := filepath.Join(homeDir, AccDataDir)

		backend := viper.GetString(FlagDBBackend)
		if backend == "" {
			backend = string(dbm.GoLevelDBBackend)
		}

		kvstore, e := types.CreateKvDB(AccSpace, types.BackendType(backend), path)
		if e != nil {
			panic("fail to open database: " + e.Error())
		}
		db := rawdb.NewDatabase(kvstore)
		gAccMptDatabase = ethstate.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     int(types.TrieCacheSize),
			Journal:   "",
			Preimages: true,
		})
	})

	return gAccMptDatabase
}

var (
	KeyPrefixLatestHeight = []byte{0x01}
	KeyPrefixRootMptHash  = []byte{0x02}
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
