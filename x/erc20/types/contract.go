package types

import (
	_ "embed"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// CompiledContract contains compiled bytecode and abi
type CompiledContract struct {
	ABI abi.ABI
	Bin string
}

var (
	EVMModuleETHAddr  common.Address
	EVMModuleBechAddr sdk.AccAddress

	// ModuleERC20Contract is the compiled okc erc20 contract
	ModuleERC20Contract CompiledContract

	//go:embed contracts/ModuleERC20.json
	moduleERC20Json []byte
)

const (
	IbcEvmModuleName = "ibc-evm"

	ContractMintMethod = "mint_by_okc_module"
)

func init() {
	EVMModuleBechAddr = authtypes.NewModuleAddress(IbcEvmModuleName)
	EVMModuleETHAddr = common.BytesToAddress(EVMModuleBechAddr.Bytes())

	if err := json.Unmarshal(moduleERC20Json, &ModuleERC20Contract); err != nil {
		panic(err)
	}
	if len(ModuleERC20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
