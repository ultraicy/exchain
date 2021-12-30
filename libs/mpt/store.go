package mpt

import (
	"io"

	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

const (
	StoreTypeMPT = types.StoreTypeIAVL
)

var (
	_ types.KVStore       = (*MptStore)(nil)
	_ types.CommitStore   = (*MptStore)(nil)
	_ types.CommitKVStore = (*MptStore)(nil)
	_ types.Queryable     = (*MptStore)(nil)
)

// MptStore Implements types.KVStore and CommitKVStore.
type MptStore struct {
	mpt Mpt
}

func (ms MptStore) GetStoreType() types.StoreType {
	return StoreTypeMPT
}

func (ms MptStore) CacheWrap() types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(ms)
}

func (ms MptStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(tracekv.NewStore(ms, w, tc))
}

func (ms MptStore) Get(key []byte) []byte {
	_, value := ms.mpt.Get(key)
	return value
}

func (ms MptStore) Has(key []byte) bool {
	return ms.mpt.Has(key)
}

func (ms MptStore) Set(key, value []byte) {
	types.AssertValidValue(value)
	ms.mpt.Set(key, value)
}

func (ms MptStore) Delete(key []byte) {
	ms.mpt.Remove(key)
}

func (ms MptStore) Iterator(start, end []byte) types.Iterator {
	//TODO implement me
	panic("implement me")
}

func (ms MptStore) ReverseIterator(start, end []byte) types.Iterator {
	//TODO implement me
	panic("implement me")
}
