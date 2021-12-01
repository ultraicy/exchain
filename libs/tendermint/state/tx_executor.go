package state

import (
	"errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"time"
)

type PreExecBlockResult struct {
	*types.Block
	*ABCIResponses
	error
}

var recordTime int64

func (blockExec *BlockExecutor) StartPreExecBlock(block *types.Block) error {
	if blockExec.processBlock != nil {
		blockExec.logger.Error("can be execute only once for a block")
		return errors.New("can be execute only once for a block")
	}

	recordTime = time.Now().UnixNano()
	blockExec.processBlock = block
	go blockExec.DoPreExecBlock(block)

	return nil
}

func (blockExec *BlockExecutor) DoPreExecBlock(block *types.Block) {
	var abciResponses *ABCIResponses
	var err error
	var preBlockRes *PreExecBlockResult
	if blockExec.isAsync {
		abciResponses, err = execBlockOnProxyAppAsync(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
	} else {
		abciResponses, err = execBlockOnProxyApp(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
	}

	if err != nil {
		preBlockRes = &PreExecBlockResult{block, abciResponses, err}
	} else {
		preBlockRes = &PreExecBlockResult{block, abciResponses, nil}
	}

	select {
	case <-blockExec.cancelChan:
		blockExec.resChan <- &PreExecBlockResult{nil, nil, errors.New("cancel_error")}
	case blockExec.resChan <- preBlockRes:
	}

	blockExec.processBlock = nil
}

func (blockExec *BlockExecutor) CancelPreExecBlock(block *types.Block) error {
	if blockExec.processBlock != block {
		blockExec.logger.Error("block: %v was cancel", block)
		return errors.New("cancel block has not begin")
	}
	// here set processBlock = nil ensure repeat call safe
	blockExec.processBlock = nil
	go func() {
		blockExec.cancelChan <- struct{}{}
	}()
	return nil
}

func (blockExec *BlockExecutor) GetPreExecBlockRes() chan *PreExecBlockResult {
	return blockExec.resChan
}

//reset base deliverState
func (blockExec *BlockExecutor) ResetDeliverState() {
	blockExec.proxyApp.SetOptionSync(abci.RequestSetOption{
		Key: "ResetDeliverState",
	})

}