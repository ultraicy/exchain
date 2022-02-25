package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

func (k Keeper) GetTokenPair(ctx sdk.Context, denom string) (types.TokenPair, error) {
	//TODO
	return types.TokenPair{}, nil
}
