package mpt

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	types3 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types2 "github.com/okex/exchain/libs/types"
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

	TriesInMemory = 100
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
	trie          ethstate.Trie
	db            ethstate.Database
	triegc        *prque.Prque
	logger        tmlog.Logger
	rootRetrieval types2.StorageRootRetrieval

	version      int64
	startVersion int64
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

func NewMptStore(logger tmlog.Logger, rootRetrieval types2.StorageRootRetrieval) *MptStore {
	db := InstanceOfAccStore()
	triegc := prque.New(nil)

	mptStore := &MptStore{
		db:            db,
		triegc:        triegc,
		logger:        logger,
		rootRetrieval: rootRetrieval,
	}
	mptStore.openTrie()

	return mptStore
}

func (ms *MptStore) openTrie() {
	latestStoredHeight := ms.GetLatestStoredBlockHeight()
	latestStoredRootHash := ms.GetMptRootHash(latestStoredHeight)

	tr, err := ms.db.OpenTrie(latestStoredRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	ms.trie = tr
	ms.version = int64(latestStoredHeight)
	ms.startVersion = int64(latestStoredHeight)

	if ms.logger != nil {
		ms.logger.Info("open acc mpt trie", "version", latestStoredHeight, "trieHash", latestStoredRootHash)
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

	root, err := ms.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		if ms.rootRetrieval != nil {
			accStorageRoot := ms.rootRetrieval(leaf)
			if accStorageRoot != types3.EmptyRootHash && accStorageRoot != (ethcmn.Hash{}) {
				ms.db.TrieDB().Reference(accStorageRoot, parent)
			}
		}
		return nil
	})
	if err != nil {
		panic("fail to commit trie data: " + err.Error())
	}
	ms.SetMptRootHash(uint64(ms.version), root)

	// TODO: use a thread to push data to database
	// push data to database
	ms.PushData2Database(ms.version)

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

func (ms *MptStore) PushData2Database(curHeight int64) {
	curMptRoot := ms.GetMptRootHash(uint64(curHeight))

	triedb := ms.db.TrieDB()
	if types2.TrieDirtyDisabled {
		if err := triedb.Commit(curMptRoot, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
		ms.SetLatestStoredBlockHeight(uint64(curHeight))
		if ms.logger != nil {
			ms.logger.Info("sync push data to db", "block", curHeight, "trieHash", curMptRoot)
		}
	} else {
		// Full but not archive node, do proper garbage collection
		triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
		ms.triegc.Push(curMptRoot, -int64(curHeight))

		if curHeight > TriesInMemory {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = ethcmn.StorageSize(256) * 1024 * 1024
			)

			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			chosen := curHeight - TriesInMemory

			if chosen <= ms.startVersion {
				return
			}

			// If the header is missing (canonical chain behind), we're reorging a low
			// diff sidechain. Suspend committing until this operation is completed.
			chRoot := ms.GetMptRootHash(uint64(chosen))
			if chRoot == (ethcmn.Hash{}) {
				if ms.logger != nil {
					ms.logger.Debug("Reorg in progress, trie commit postponed", "number", chosen)
				}
			} else {
				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				if err := triedb.Commit(chRoot, true, nil); err != nil {
					panic("fail to commit mpt data: " + err.Error())
				}
				ms.SetLatestStoredBlockHeight(uint64(chosen))
				if ms.logger != nil {
					ms.logger.Info("async push data to db", "block", chosen, "trieHash", chRoot)
				}
			}

			// Garbage collect anything below our required write retention
			for !ms.triegc.Empty() {
				root, number := ms.triegc.Pop()
				if int64(-number) > chosen {
					ms.triegc.Push(root, number)
					break
				}
				triedb.Dereference(root.(ethcmn.Hash))
			}
		}
	}
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (ms *MptStore) OnStop() error {
	// Ensure the state of a recent block is also stored to disk before exiting.
	// We're writing three different states to catch different restart scenarios:
	//  - HEAD:     So we don't need to reprocess any blocks in the general case
	//  - HEAD-1:   So we don't do large reorgs if our HEAD becomes an uncle
	//  - HEAD-127: So we have a hard limit on the number of blocks reexecuted
	if !types2.TrieDirtyDisabled {
		triedb := ms.db.TrieDB()
		oecStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of oec

		for _, offset := range []uint64{0, 1, TriesInMemory - 1} {
			if number := uint64(ms.version); number > offset {
				recent := number - offset
				if recent <= oecStartHeight || recent <= uint64(ms.startVersion) {
					break
				}

				recentMptRoot := ms.GetMptRootHash(recent)
				if recentMptRoot == (ethcmn.Hash{}) {
					if ms.logger != nil {
						ms.logger.Debug("Reorg in progress, trie commit postponed", "block", recent)
					}
				} else {
					if ms.logger != nil {
						ms.logger.Info("Writing cached state to disk", "block", recent, "trieHash", recentMptRoot)
					}
					if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
						if ms.logger != nil {
							ms.logger.Error("Failed to commit recent state trie", "err", err)
						}
					}
				}
			}
		}

		for !ms.triegc.Empty() {
			ms.db.TrieDB().Dereference(ms.triegc.PopItem().(ethcmn.Hash))
		}
	}

	return nil
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
	trie, err := ms.getHeight(req)
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
func (ms *MptStore) getHeight(req abci.RequestQuery) (ethstate.Trie, error) {
	height := uint64(req.Height)
	latestStoredBlockHeight := ms.GetLatestStoredBlockHeight()
	if height == 0 || height > latestStoredBlockHeight {
		height = latestStoredBlockHeight
	}

	latestRootHash := ms.GetMptRootHash(height)
	return ms.db.OpenTrie(latestRootHash)
}

func getVersioned(trie ethstate.Trie, key []byte) ([]byte, error) {
	return trie.TryGet(key)
}

// getVersionedWithProof returns the Merkle proof for given storage slot.
func getVersionedWithProof(trie ethstate.Trie, key []byte) ([]byte, [][]byte, error) {
	value, err := trie.TryGet(key)
	if err != nil {
		return nil, nil, err
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
