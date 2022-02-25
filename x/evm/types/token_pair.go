package types

type ContractOwner int32

const (
	// OwnerModule erc20 is owned by the erc20 module account.
	// first native coin, then erc20 coin
	OwnerModule ContractOwner = iota
	// OwnerExternal erc20 is owned by an external account.
	// first erc20 coin, then native coin
	OwnerExternal
)

type TokenPair struct {
	// address of ERC20 contract token
	ERC20ContractAddress string `json:"erc20_contract_address"`
	// cosmos base denomination to be mapped to
	Denom string `json:"denom"`
	// ERC20 owner address (0 ModuleAccount, 1 external address)
	ContractOwner `json:"contract_owner"`
}
