package rootmulti

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func queryIbcProof(res *abci.ResponseQuery, info *commitInfo, storeName string) {
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}

func (s *Store) getStores() map[types.StoreKey]types.CommitKVStore {
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
