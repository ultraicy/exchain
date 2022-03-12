package rootmulti

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

func queryIbcProof(res *abci.ResponseQuery, info *commitInfo, storeName string) {
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}

type StoreOption func(s *Store)

func WithHeightFilterPipeline(f HeightFilterPipeline) StoreOption {
	return func(s *Store) {
		s.commitHeightFilterPipeline = merge(f, s.commitHeightFilterPipeline)
	}
}
func WithPruneHeightBlockFilter(f HeightFilterPipeline) StoreOption {
	return func(s *Store) {
		s.pruneHeightFilterPipeline = merge(f, s.pruneHeightFilterPipeline)
	}
}
func merge(f, s HeightFilterPipeline) HeightFilterPipeline {
	return func(h int64) func(str string) bool {
		filter := f(h)
		if nil != filter {
			return filter
		}
		return s(h)
	}
}

func (s *Store) getStores(h int64) map[types.StoreKey]types.CommitKVStore {
	return s.stores
}

func (s *Store) getFilterStores(h int64) map[types.StoreKey]types.CommitKVStore {
	if tmtypes.HigherThanIBCHeight(h) {
		return s.stores
	}
	f := s.pruneHeightFilterPipeline(h)
	// TODO FILTER:
	m := make(map[types.StoreKey]types.CommitKVStore)
	for k, v := range s.stores {
		if f(k.Name()) {
			continue
		}
		m[k] = v
	}
	return m
}
