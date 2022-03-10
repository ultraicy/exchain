package keeper

import (
	"fmt"
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	"sync/atomic"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical Account constructor.
	proto func() exported.Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace subspace.Subspace

	//permAddrs map[string]types.PermissionsForAddress

	observers []ObserverI

	count int32
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
) AccountKeeper {
	return NewAccountKeeperWithPer(cdc, key, paramstore, proto, nil)
}

func NewAccountKeeperWithPer(
	cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
	maccPerms map[string][]string,
) AccountKeeper {
	//// set KeyTable if it has not already been set
	//if !paramstore.HasKeyTable() {
	//	paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	//}
	//
	//permAddrs := make(map[string]types.PermissionsForAddress)
	//if maccPerms != nil {
	//	for name, perms := range maccPerms {
	//		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	//	}
	//}
	//
	//return AccountKeeper{
	//	key:           key,
	//	proto:         proto,
	//	cdc:           cdc,
	//	paramSubspace: paramstore,
	//	permAddrs:     permAddrs,
	//}

	return AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (ak AccountKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, error) {
	logrusplugin.Info("count", "count", atomic.AddInt32(&ak.count, 1))
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error) {
	logrusplugin.Info("count", "count", atomic.AddInt32(&ak.count, 1),"add",addr.String())
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetSequence(), nil
}

// GetNextAccountNumber returns and increments the global account number counter.
// If the global account number is not set, it initializes it with value 0.
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	logrusplugin.Info("count", "count", atomic.AddInt32(&ak.count, 1))
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		// initialize the account numbers
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(types.GlobalAccountNumberKey, bz)

	return accNumber
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	logrusplugin.Info("count", "count", atomic.AddInt32(&ak.count, 1))
	val, err := ak.cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(bz, &acc)
	if err == nil {
		acc = val.(exported.Account)
		return
	}
	err = ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}

// GetModuleAddress returns an address based on the module name
//func (ak AccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
//	permAddr, ok := ak.permAddrs[moduleName]
//	if !ok {
//		return nil
//	}
//
//	return permAddr.GetAddress()
//}

// GetModuleAccount gets the module account from the auth account store, if the account does not
// exist in the AccountKeeper, then it is created.
//func (ak AccountKeeper) GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI {
//	acc, _ := ak.GetModuleAccountAndPermissions(ctx, moduleName)
//	return acc
//}

// GetModuleAccountAndPermissions gets the module account from the auth account store and its
// registered permissions
//func (ak AccountKeeper) GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string) {
//	addr, perms := ak.GetModuleAddressAndPermissions(moduleName)
//	if addr == nil {
//		return nil, []string{}
//	}
//
//	acc := ak.GetAccount(ctx, addr)
//	if acc != nil {
//		macc, ok := acc.(types.ModuleAccountI)
//		if !ok {
//			panic("account is not a module account")
//		}
//		return macc, perms
//	}
//
//	// create a new module account
//	macc := types.NewEmptyModuleAccount(moduleName, perms...)
//	maccI := (ak.NewAccount(ctx, macc)).(types.ModuleAccountI) // set the account number
//	ak.SetModuleAccount(ctx, maccI)
//
//	return maccI, perms
//}
