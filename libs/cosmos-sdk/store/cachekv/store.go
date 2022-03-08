package cachekv

import (
	"bytes"
	"container/list"
	"io"
	"reflect"
	"sort"
	"sync"
	"unsafe"

	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store struct {
	mtx           sync.Mutex
	dirtyCache    map[string]types.CValue
	unsortedCache map[string]struct{}
	sortedCache   *list.List // always ascending sorted
	parent        types.KVStore
	readList      map[string][]byte
}

var _ types.CacheKVStore = (*Store)(nil)

func NewStore(parent types.KVStore) *Store {
	return &Store{
		dirtyCache:    make(map[string]types.CValue),
		readList:      make(map[string][]byte),
		unsortedCache: make(map[string]struct{}),
		sortedCache:   list.New(),
		parent:        parent,
	}
}

// Implements Store.
func (store *Store) GetStoreType() types.StoreType {
	return store.parent.GetStoreType()
}

// Implements types.KVStore.
func (store *Store) Get(key []byte) (value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)

	cacheValue, ok := store.dirtyCache[string(key)]
	if !ok {
		value = store.parent.Get(key)
		store.setCacheValue(key, value, false, false)
	} else {
		value = cacheValue.Value
	}

	return value
}

func (store *Store) IteratorCache(cb func(key, value []byte, isDirty bool, isDelete bool, sKey types.StoreKey) bool, sKey types.StoreKey) bool {
	if cb == nil || len(store.dirtyCache) == 0 {
		return true
	}
	store.mtx.Lock()
	defer store.mtx.Unlock()

	for key, v := range store.dirtyCache {
		if !cb([]byte(key), v.Value, v.Dirty, v.Deleted, sKey) {
			return false
		}
	}
	return true
}

func (store *Store) GetInitRead() (map[types.StoreKey]map[string][]byte, map[types.StoreKey]map[string]types.CValue) {
	readList := make(map[types.StoreKey]map[string][]byte)
	readList[types.NullKey] = make(map[string][]byte)
	for k, v := range store.readList {
		readList[types.NullKey][k] = v
	}

	dirtyList := make(map[types.StoreKey]map[string]types.CValue)
	dirtyList[types.NullKey] = make(map[string]types.CValue)
	for k, v := range store.dirtyCache {
		dirtyList[types.NullKey][k] = v
	}
	return readList, dirtyList
}

// Implements types.KVStore.
func (store *Store) Set(key []byte, value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)
	types.AssertValidValue(value)

	store.setCacheValue(key, value, false, true)
}

// Implements types.KVStore.
func (store *Store) Has(key []byte) bool {
	value := store.Get(key)
	return value != nil
}

// Implements types.KVStore.
func (store *Store) Delete(key []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)

	store.setCacheValue(key, nil, true, true)
}

// Implements Cachetypes.KVStore.
func (store *Store) Write() {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, 0, len(store.dirtyCache))
	for key, _ := range store.dirtyCache {
		keys = append(keys, key)

	}

	sort.Strings(keys)

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		cacheValue := store.dirtyCache[key]
		switch {
		case cacheValue.Deleted:
			store.parent.Delete([]byte(key))
		case cacheValue.Value == nil:
			// Skip, it already doesn't exist in parent.
		default:
			store.parent.Set([]byte(key), cacheValue.Value)
		}
	}

	// Clear the cache
	store.dirtyCache = make(map[string]types.CValue)
	store.unsortedCache = make(map[string]struct{})
	store.sortedCache.Init()
}

//----------------------------------------
// To cache-wrap this Store further.

// Implements CacheWrapper.
func (store *Store) CacheWrap() types.CacheWrap {
	return NewStore(store)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (store *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return NewStore(tracekv.NewStore(store, w, tc))
}

//----------------------------------------
// Iteration

// Implements types.KVStore.
func (store *Store) Iterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, true)
}

// Implements types.KVStore.
func (store *Store) ReverseIterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, false)
}

func (store *Store) iterator(start, end []byte, ascending bool) types.Iterator {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	var parent, cache types.Iterator

	if ascending {
		parent = store.parent.Iterator(start, end)
	} else {
		parent = store.parent.ReverseIterator(start, end)
	}

	store.dirtyItems(start, end)
	cache = newMemIterator(start, end, store.sortedCache, ascending)

	return newCacheMergeIterator(parent, cache, ascending)
}

// strToByte is meant to make a zero allocation conversion
// from string -> []byte to speed up operations, it is not meant
// to be used generally, but for a specific pattern to check for available
// keys within a domain.
func strToByte(s string) []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = len(s)
	hdr.Len = len(s)
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	return b
}

// byteSliceToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func byteSliceToStr(b []byte) string {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(hdr))
}

// Constructs a slice of dirty items, to use w/ memIterator.
func (store *Store) dirtyItems(start, end []byte) {
	unsorted := make([]*tmkv.Pair, 0)

	n := len(store.unsortedCache)
	for key := range store.unsortedCache {
		if dbm.IsKeyInDomain(strToByte(key), start, end) {
			cacheValue := store.dirtyCache[key]
			unsorted = append(unsorted, &tmkv.Pair{Key: []byte(key), Value: cacheValue.Value})
		}
	}

	if len(unsorted) == n { // This pattern allows the Go compiler to emit the map clearing idiom for the entire map.
		for key := range store.unsortedCache {
			delete(store.unsortedCache, key)
		}
	} else { // Otherwise, normally delete the unsorted keys from the map.
		for _, kv := range unsorted {
			delete(store.unsortedCache, byteSliceToStr(kv.Key))
		}
	}

	sort.Slice(unsorted, func(i, j int) bool {
		return bytes.Compare(unsorted[i].Key, unsorted[j].Key) < 0
	})

	for e := store.sortedCache.Front(); e != nil && len(unsorted) != 0; {
		uitem := unsorted[0]
		sitem := e.Value.(*tmkv.Pair)
		comp := bytes.Compare(uitem.Key, sitem.Key)
		switch comp {
		case -1:
			unsorted = unsorted[1:]
			store.sortedCache.InsertBefore(uitem, e)
		case 1:
			e = e.Next()
		case 0:
			unsorted = unsorted[1:]
			e.Value = uitem
			e = e.Next()
		}
	}

	for _, kvp := range unsorted {
		store.sortedCache.PushBack(kvp)
	}

}

//----------------------------------------
// etc

// Only entrypoint to mutate store.cache.
func (store *Store) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	keyStr := string(key)
	if !dirty {
		store.readList[keyStr] = value
		return
	}

	store.dirtyCache[keyStr] = types.CValue{
		Value:   value,
		Deleted: deleted,
		Dirty:   dirty,
	}
	if dirty {
		store.unsortedCache[keyStr] = struct{}{}
	}
}
