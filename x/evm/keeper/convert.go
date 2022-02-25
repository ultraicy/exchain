package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// ConvertCoin4NativeCoin handles the Coin conversion flow for a native coin
// token pair:
//  - Bank: Escrow Coins on module account (Coins are not burned)
//  - Contract: Mint Tokens and send to receiver
//  - Check if token balance increased by amount
func (k Keeper) ConvertCoin4NativeCoin(
	ctx sdk.Context,
	pair types.TokenPair,
	msg *types.MsgConvertCoin,
	receiver common.Address,
	sender sdk.AccAddress) (*sdk.Result, error) {
	return nil, nil
}

// ConvertCoin4NativeERC20 handles the Coin conversion flow for a native ERC20
// token pair:
//  - Bank: Escrow Coins on module account
//  - Contract: Unescrow Tokens that have been previously escrowed with ConvertERC20 and send to receiver
//  - Bank: Burn escrowed Coins
//  - Check if token balance increased by amount
//  - Check for unexpected `appove` event in logs
func (k Keeper) ConvertCoin4NativeERC20(
	ctx sdk.Context,
	pair types.TokenPair,
	msg *types.MsgConvertCoin,
	receiver common.Address,
	sender sdk.AccAddress) (*sdk.Result, error) {
	return nil, nil
}
