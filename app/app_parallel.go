package app

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	authTypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"sync"
)

// feeCollectorHandler set or get the value of feeCollectorAcc
func updateFeeCollectorHandler(bk bank.Keeper, sk supply.Keeper) sdk.UpdateFeeCollectorAccHandler {
	return func(ctx sdk.Context, balance sdk.Coins) error {
		return bk.SetCoins(ctx, sk.GetModuleAddress(auth.FeeCollectorName), balance)
	}
}

// evmTxFeeHandler get tx fee for evm tx
func evmTxFeeHandler() sdk.GetTxFeeHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool, signCache sdk.SigCache) {
		if evmTx, ok := tx.(evmtypes.MsgEthereumTx); ok {
			isEvm = true
			signCache, _ = evmTx.VerifySig(evmTx.ChainID(), ctx.BlockHeight(), nil)

		}
		if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
		}

		return
	}
}

// fixLogForParallelTxHandler fix log for parallel tx
func fixLogForParallelTxHandler(ek *evm.Keeper) sdk.LogFix {
	return func(execResults [][]string) (logs [][]byte) {
		return ek.FixLog(execResults)
	}
}

func preLoadSender(ak auth.AccountKeeper, key sdk.StoreKey) sdk.PreLoadSender {
	return func(ctx sdk.Context, addr sdk.AccAddress, mu *sync.Mutex) {
		return

		store := ctx.KVStore(key)
		authAddr := authTypes.AddressStoreKey(addr)
		bz := store.Get(authAddr)
		if bz == nil {
			mu.Lock()
			ctx.Cache().UpdateAccount(authAddr, nil, len(bz), false)
			mu.Unlock()
			return
		}
		acc := ak.DecodeAccount(bz)
		mu.Lock()
		ctx.Cache().UpdateAccount(authAddr, acc, len(bz), false)
		mu.Unlock()

	}
}
