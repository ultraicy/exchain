package mpt

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
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

var cdc = codec.New()

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
	logger tmlog.Logger

	version int64
}

func (ms *MptStore) GetFlatKVReadTime() int {
	return 0
}

func (ms *MptStore) GetFlatKVWriteTime() int {
	return 0
}

func (ms *MptStore) GetFlatKVReadCount() int {
	return 0
}

func (ms *MptStore) GetFlatKVWriteCount() int {
	return 0
}

func NewMptStore(logger tmlog.Logger) *MptStore {
	db := InstanceOfAccStore()
	triegc := prque.New(nil)

	latestHeight := GetLatestStoredBlockHeight(db)
	if logger != nil {
		logger.Info("latest stored Block", "height", latestHeight)
	}

	latestRootHash := GetMptRootHash(db, latestHeight)
	if logger != nil {
		logger.Info("latest mpt hash", "hash", latestRootHash)
	}

	tr, err := db.OpenTrie(latestRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}

	return &MptStore{
		tr,
		db,
		triegc,
		logger,
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
func (ms *MptStore) Commit(delta *iavl.TreeDelta, bytes []byte) (types.CommitID, iavl.TreeDelta, []byte) {
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
	ms.db.TrieDB().DiskDB().Put(KeyPrefixLatestHeight, height)

	if ms.logger != nil {
		ms.logger.Info("acc mpt commit", "height", ms.version, "hash", root.String())
	}

	return types.CommitID{
		Version: ms.version,
		Hash:    root.Bytes(),
	}, iavl.TreeDelta{}, nil
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
func (ms *MptStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrTxDecode, "query cannot be zero length"))
	}

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	trie, err := getHeight(ms.db, req)
	if err != nil {
		res.Log = iavl.ErrVersionDoesNotExist.Error()
		return
	}

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if req.Prove {
			value, proof, err := getVersionedWithProof(trie, key)
			if err != nil {
				res.Log = err.Error()
				break
			}
			if proof == nil {
				// Proof == nil implies that the store is empty.
				if value != nil {
					panic("unexpected value for an empty proof")
				}
			}
			if value != nil {
				// value was found
				res.Value = value
				//TODO: translate proof to RangeProof
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewValueOp(key, nil).ProofOp()}}
			} else {
				// value wasn't found
				res.Value = nil
				//TODO: translate proof to RangeProof
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewAbsenceOp(key, nil).ProofOp()}}
			}
		} else {
			res.Value, _ = getVersioned(trie, key)
		}

	case "/subspace":
		var KVs []types.KVPair

		subspace := req.Data
		res.Key = subspace

		iterator := newMptIterator(trie, subspace, sdk.PrefixEndBytes(subspace))
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, types.KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value = cdc.MustMarshalBinaryLengthPrefixed(KVs)

	default:
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unexpected query path: %v", req.Path))
	}

	return res
}

// Handle gatest the latest height, if height is 0
func getHeight(db ethstate.Database, req abci.RequestQuery) (ethstate.Trie, error) {
	height := uint64(req.Height)
	latestStoredBlockHeight := GetLatestStoredBlockHeight(db)
	if height == 0 || height > latestStoredBlockHeight{
		height = latestStoredBlockHeight
	}

	latestRootHash := GetMptRootHash(db, height)
	return db.OpenTrie(latestRootHash)
}

func getVersioned(trie ethstate.Trie, key []byte) ([]byte, error) {
	return trie.TryGet(key)
}

// getVersionedWithProof returns the Merkle proof for given storage slot.
func getVersionedWithProof(trie ethstate.Trie, key []byte) ([]byte, [][]byte, error) {
	value, err := trie.TryGet(key)
	if err != nil {
		return nil, nil ,err
	}

	var proof proofList
	err = trie.Prove(crypto.Keccak256(key), 0, &proof)
	return value, proof, err
}

type proofList [][]byte

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}