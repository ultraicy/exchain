package mpt

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/prque"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

const (
	StoreTypeMPT = types.StoreTypeMPT
)

var (
	_ types.KVStore       = (*MptStore)(nil)
	_ types.CommitStore   = (*MptStore)(nil)
	_ types.CommitKVStore = (*MptStore)(nil)
	_ types.Queryable     = (*MptStore)(nil)
)

// MptStore Implements types.KVStore and CommitKVStore.
type MptStore struct {
	trie   ethstate.Trie
	db     ethstate.Database
	triegc *prque.Prque

	version int64
}

func NewMptStore() *MptStore {
	db := InstanceOfEvmStore()
	triegc := prque.New(nil)

	latestHeight := GetLatestBlockHeight(db)
	fmt.Println("latest Block height", latestHeight)
	latestRootHash := GetRootMptHash(db, latestHeight)
	fmt.Println("latest MPT hash", latestRootHash)
	tr, err := db.OpenTrie(latestRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}

	return &MptStore{
		tr,
		db,
		triegc,
		int64(latestHeight),
	}
}

/*
*  implement KVStore
 */
func (ms *MptStore) GetStoreType() types.StoreType {
	return StoreTypeMPT
}

func (ms *MptStore) CacheWrap() types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(ms)
}

func (ms *MptStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(tracekv.NewStore(ms, w, tc))
}

func (ms *MptStore) Get(key []byte) []byte {
	value, err := ms.trie.TryGet(key)
	if err != nil {
		return nil
	}
	return value
}

func (ms *MptStore) Has(key []byte) bool {
	return ms.Get(key) != nil
}

func (ms *MptStore) Set(key, value []byte) {
	types.AssertValidValue(value)

	err := ms.trie.TryUpdate(key, value)
	if err != nil {
		return
	}
	return
}

func (ms *MptStore) Delete(key []byte) {
	err := ms.trie.TryDelete(key)
	if err != nil {
		return
	}
}

func (ms *MptStore) Iterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.trie, start, end)
}

func (ms *MptStore) ReverseIterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.trie, start, end)
}

/*
*  implement CommitStore, CommitKVStore
 */
func (ms *MptStore) Commit() types.CommitID {
	ms.version++
	root, err := ms.trie.Commit(nil)
	if err != nil {
		panic("fail to commit trie data: " + err.Error())
	}
	err = ms.db.TrieDB().Commit(root, true, nil)
	if err != nil {
		panic("fail to commit trieDB data: " + err.Error())
	}

	height := sdk.Uint64ToBigEndian(uint64(ms.version))
	ms.db.TrieDB().DiskDB().Put(append(KeyPrefixRootMptHash, height...), root.Bytes())
	fmt.Println("set root hash", root.String())
	ms.db.TrieDB().DiskDB().Put(KeyPrefixLatestHeight, height)
	fmt.Println("set Block height", ms.version)

	return types.CommitID{
		Version: ms.version,
		Hash:    root.Bytes(),
	}
}

func (ms *MptStore) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: ms.version,
		Hash:    ms.trie.Hash().Bytes(),
	}
}

func (ms *MptStore) SetPruning(options types.PruningOptions) {
	panic("cannot set pruning options on an initialized MPT store")
}

func (ms *MptStore) GetDBWriteCount() int {
	return 0
}

func (ms *MptStore) GetDBReadCount() int {
	return 0
}

func (ms *MptStore) GetNodeReadCount() int {
	return 0
}

func (ms *MptStore) ResetCount() {
	return
}

func (ms *MptStore) GetDBReadTime() int {
	return 0
}

/*
*  implement Queryable
 */
func (ms *MptStore) Query(query abci.RequestQuery) abci.ResponseQuery {
	//TODO implement me
	panic("implement me")
}
