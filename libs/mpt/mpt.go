package mpt

import (
	"sync/atomic"

	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	iavltree "github.com/okex/exchain/libs/iavl"
)

var _ iavl.Tree = (*Mpt)(nil)

type Mpt struct {
	name string

	readCount  int64
	writeCount int64

	trie ethstate.Trie
	db   ethstate.Database
}

func (m Mpt) Hash() []byte {
	return m.trie.Hash().Bytes()
}

func (m Mpt) Has(key []byte) bool {
	_, value := m.Get(key)
	return value != nil
}

func (m Mpt) Get(key []byte) (index int64, value []byte) {
	enc, err := m.trie.TryGet(key)
	if err != nil {
		return 0, nil
	}
	if len(enc) == 0 {
		return 0, nil
	}
	return 0, enc
}

func (m Mpt) Set(key, value []byte) bool {
	err := m.trie.TryUpdate(key, value)
	if err != nil {
		return false
	}
	return true
}

func (m Mpt) Remove(key []byte) ([]byte, bool) {
	err := m.trie.TryDelete(key)
	if err != nil {
		return nil, false
	}
	return nil, true
}
func (m Mpt) GetModuleName() string {
	return m.name
}

func (m Mpt) GetDBWriteCount() int {
	return int(atomic.LoadInt64(&m.writeCount))
}

func (m Mpt) GetDBReadCount() int {
	return int(atomic.LoadInt64(&m.readCount))
}

func (m Mpt) GetNodeReadCount() int {
	return 0
}

func (m Mpt) ResetCount() {
	atomic.StoreInt64(&m.writeCount, 0)
	atomic.StoreInt64(&m.readCount, 0)
}

func (m Mpt) SaveVersion() ([]byte, int64, error) {
	panic("should not be executed in mpt")
}

func (m Mpt) DeleteVersion(version int64) error {
	panic("should not be executed in mpt")
}

func (m Mpt) DeleteVersions(versions ...int64) error {
	panic("should not be executed in mpt")
}

func (m Mpt) Version() int64 {
	return 0
}

func (m Mpt) VersionExists(version int64) bool {
	//TODO implement me
	panic("implement me")
}

func (m Mpt) GetVersioned(key []byte, version int64) (int64, []byte) {
	//TODO implement me
	panic("implement me")
}

func (m Mpt) GetVersionedWithProof(key []byte, version int64) ([]byte, *iavltree.RangeProof, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mpt) GetImmutable(version int64) (*iavltree.ImmutableTree, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mpt) SetInitialVersion(version uint64) {
	//TODO implement me
	panic("implement me")
}
