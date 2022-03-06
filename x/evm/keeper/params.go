package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if ctx.BlockHeight() > types2.GetAnteHeight() && ctx.IsDeliverorAsync() {
		if types.EvmParamsCache.IsNeedParamsUpdate() {
			k.paramSpace.GetParamSet(ctx, &params)
			types.EvmParamsCache.UpdateParams(params)
		} else {
			params = types.EvmParamsCache.GetParams()
		}
	} else {
		k.paramSpace.GetParamSet(ctx, &params)
	}

	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	types.EvmParamsCache.SetNeedParamsUpdate()
}
