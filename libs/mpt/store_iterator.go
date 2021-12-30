package mpt

import (
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

var _ types.Iterator = (*mptIterator)(nil)

type mptIterator struct {
	// Domain
	start, end []byte

	*trie.Iterator
}

func newMptIterator(mpt *Mpt, start, end []byte, ascending bool) *mptIterator {
	iter := &mptIterator{
		start:    start,
		end:      end,
		Iterator: trie.NewIterator(mpt.trie.NodeIterator(start)),
	}
	return iter
}

func (it mptIterator) Domain() (start []byte, end []byte) {
	return it.start, it.end
}

func (it mptIterator) Valid() bool {
	// return it.invalid
	return true
}

func (it mptIterator) Next() {
	it.Iterator.Next()
}

func (it mptIterator) Key() (key []byte) {
	return it.Iterator.Key
}

func (it mptIterator) Value() (value []byte) {
	return it.Iterator.Value
}

func (it mptIterator) Error() error {
	return it.Iterator.Err
}

func (it mptIterator) Close() {
	return
}
