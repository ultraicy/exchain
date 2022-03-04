package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"math/big"
	"sync"

	"github.com/okex/exchain/x/evm/types"
)

func (k *Keeper) FixLog(txsInfo []*sdk.ParaTxInfo) [][]byte {

	txSize := len(txsInfo)

	res := make([][]byte, txSize, txSize)
	logSize := uint(0)
	txInBlock := -1
	k.Bloom = new(big.Int)

	for index, info := range txsInfo {
		if info.AnteErr != nil {
			//fmt.Println("zzzzz-1", index, info.AnteErr)
			continue
		}

		rs, ok := k.LogsManages.Results[info.ResultID]
		if !ok {
			//fmt.Println("zzzzz-2", index, ok, info.ResultID)
			continue
		}

		txInBlock++
		if rs == nil {
			//fmt.Println("zzzzz-3", index)
			continue
		}

		for _, v := range rs.Logs {
			v.Index = logSize
			v.TxIndex = uint(txInBlock)
			logSize++
		}

		//fmt.Println("bloom", index, ethcommon.BytesToHash(rs.Bloom.Bytes()).String())
		k.Bloom = k.Bloom.Or(k.Bloom, rs.Bloom.Big())
		data, err := types.EncodeResultData(*rs)
		if err != nil {
			panic(err)
		}
		res[index] = data
	}
	return res
}

type LogsManager struct {
	cnt int

	mu         sync.RWMutex
	txMapIndex map[string]int
	Results    map[int]*types.ResultData
}

func NewLogManager() *LogsManager {
	return &LogsManager{
		mu:         sync.RWMutex{},
		txMapIndex: make(map[string]int, 0),
		Results:    make(map[int]*types.ResultData),
	}
}

func (l *LogsManager) Set(txBytes string, value *types.ResultData) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.txMapIndex[txBytes] = l.cnt
	l.Results[l.cnt] = value
	l.cnt++
}

func (l *LogsManager) GetResultID(txBytes string) (int, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	data, ok := l.txMapIndex[txBytes]
	delete(l.txMapIndex, txBytes)
	return data, ok
}

func (l *LogsManager) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Results)
}

func (l *LogsManager) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.txMapIndex = make(map[string]int)
	l.Results = make(map[int]*types.ResultData)
}
