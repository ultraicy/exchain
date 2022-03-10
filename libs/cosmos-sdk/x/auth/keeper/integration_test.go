package keeper_test

import (
	"encoding/hex"
	"fmt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}

func TestADD(t *testing.T) {
	str := "010dfe79af583628e13ecfd34d021f1d1e99fbf9208ef66c9f27e4a2c41f2d95"
	vs, err := hex.DecodeString(str)
	if nil != err {
		panic(err)
	}
	fromHex, err := sdk.AccAddressFromHex(str)
	fmt.Println(err)
	ret := fromHex.Bech32String("ex")
	fmt.Println(ret)
	fmt.Println(string(vs))
}
