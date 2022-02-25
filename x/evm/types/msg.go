package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgConvertCoin{}
	_ sdk.Msg = &MsgConvertERC20{}
)

const (
	TypeMsgConvertCoin  = "convert_coin"
	TypeMsgConvertERC20 = "convert_ERC20"
)

// MsgConvertCoin defines a Msg to convert a oec Coin to a ERC20 token
type MsgConvertCoin struct {
	// oec coin which denomination is registered on erc20 bridge.
	// The coin amount defines the total ERC20 tokens to convert.
	Coin types.SysCoin `json:"coin"`
	// recipient hex address to receive ERC20 token
	Receiver string `json:"receiver"`
	// oec bech32 address from the owner of the given ERC20 tokens
	Sender string `json:"sender"`
}

func (msg MsgConvertCoin) Route() string {
	return RouterKey
}

func (msg MsgConvertCoin) Type() string {
	return TypeMsgConvertCoin
}

func (msg MsgConvertCoin) ValidateBasic() error {
	if err := ValidateErc20Denom(msg.Coin.Denom); err != nil {
		return err
	}

	if !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cannot mint a non-positive amount")
	}
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	if !common.IsHexAddress(msg.Receiver) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver hex address %s", msg.Receiver)
	}
	return nil
}

func (msg MsgConvertCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgConvertCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateErc20Denom checks if a denom is a valid erc20 denomination
func ValidateErc20Denom(denom string) error {
	denomSplit := strings.SplitN(denom, "/", 2)

	if len(denomSplit) != 2 || denomSplit[0] != ModuleName {
		return fmt.Errorf("invalid denom. %s denomination should be prefixed with the format 'erc20/", denom)
	}

	if !common.IsHexAddress(denomSplit[1]) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress, "address '%s' is not a valid ethereum hex address",
			denomSplit[1],
		)
	}
	return nil
}

// MsgConvertERC20 defines a Msg to convert an ERC20 token to a oec coin.
type MsgConvertERC20 struct {
	// ERC20 token contract address registered on erc20 bridge
	ContractAddress string `json:"contract_address"`
	// amount of ERC20 tokens to mint
	Amount sdk.Int `json:"amount"`
	// bech32 address to receive SDK coins.
	Receiver string `json:"receiver"`
	// sender hex address from the owner of the given ERC20 tokens
	Sender string `json:"sender"`
}

func (msg MsgConvertERC20) Route() string { return RouterKey }

func (msg MsgConvertERC20) Type() string { return TypeMsgConvertERC20 }

func (msg MsgConvertERC20) ValidateBasic() error {
	if !common.IsHexAddress(msg.ContractAddress) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract hex address '%s'", msg.ContractAddress)
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cannot mint a non-positive amount")
	}
	_, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid reciver address")
	}
	if !common.IsHexAddress(msg.Sender) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender hex address %s", msg.Sender)
	}
	return nil
}

func (msg MsgConvertERC20) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgConvertERC20) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{common.HexToAddress(msg.Sender).Bytes()}
}
