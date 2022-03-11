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
		s.heightFilterPipeline = merge(f, s.heightFilterPipeline)
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
	// TODO FILTER:
	m := make(map[types.StoreKey]types.CommitKVStore)
	b := make(map[string]struct{})
	b["ibc"] = struct{}{}
	b["mem_capability"] = struct{}{}
	b["capability"] = struct{}{}
	b["transfer"] = struct{}{}
	for k, v := range s.stores {
		if _, exist := b[k.Name()]; exist {
			continue
		}
		m[k] = v
	}
	return m
}
