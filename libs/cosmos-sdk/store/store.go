package store

import (
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/cache"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB, ops ...rootmulti.StoreOption) types.CommitMultiStore {
	return rootmulti.NewStore(db, ops...)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
