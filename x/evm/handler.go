package evm

import (
	"github.com/okex/exchain/x/evm/txs"
	"github.com/okex/exchain/x/evm/txs/base"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	common2 "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/evm/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		evmtx, ok := msg.(types.MsgEthereumTx)
		if ok {
			result, err = handleMsgEthereumTx(ctx, k, &evmtx)
			if err != nil {
				err = sdkerrors.New(types.ModuleName, types.CodeSpaceEvmCallFailed, err.Error())
			}
		} else {
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		return result, err
	}
}


func hgu(gc int64, tx *types.MsgEthereumTx) (err error) {
	if cfg.DynamicConfig.GetMaxGasUsedPerBlock() < 0 {
		return
	}

	db := bam.InstanceOfHistoryGasUsedRecordDB()
	msgFnSignature, toDeployContractSize := tx.GetTxFnSignatureInfo()

	if msgFnSignature == nil {
		return
	}

	var hisGu []byte
	hisGu, err = db.Get(msgFnSignature)
	if err != nil {
		return
	}

	if toDeployContractSize > 0 {
		// calculate average gas consume for deploy contract case
		gc = gc / int64(toDeployContractSize)
	}

	var avgGas int64
	if hisGu != nil {
		hgu := common2.BytesToInt64(hisGu)
		avgGas = int64(bam.GasUsedFactor*float64(gc) + (1.0-bam.GasUsedFactor)*float64(hgu))
	} else {
		avgGas = gc
	}

	return db.Set(msgFnSignature, common2.Int64ToBytes(avgGas))
}

func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg *types.MsgEthereumTx) (*sdk.Result, error) {
	txFactory := txs.NewFactory(base.Config{
		Ctx:    ctx,
		Keeper: k,
	})
	var res *sdk.Result
	tx, err := txFactory.CreateTx()
	if err == nil {
		res, err = txs.TransitionEvmTx(tx, msg)
	}

	if err == nil {
		hgu(int64(ctx.GasMeter().GasConsumed()), msg)
	}

	return res, err
}



