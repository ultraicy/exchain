package types

//import (
//	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
//	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
//	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
//	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
//	"strings"
//)

//
//import (
//	"github.com/gogo/protobuf/proto"
//	"github.com/okex/exchain/libs/cosmos-sdk/types"
//)

//var (
//	_ types.MsgAdapter = (*MsgTransfer)(nil)
//)
//
//func (m *MsgTransfer) Reset()         {}
//func (m *MsgTransfer) String() string { return proto.CompactTextString(m)  }
//func (m *MsgTransfer) ProtoMessage()  {}
//func(m *MsgTransfer)ToMsgTransfer()*MsgTransfer{
//	ret:=&MsgTransfer{
//		SourcePort:       m.SourcePort,
//		SourceChannel:    m.SourceChannel,
//		Token:            sdk.Coin{
//			Denom:  m.Token.Denom,
//			Amount: sdk.Dec{
//				Int: m.Token.Amount.BigInt(),
//			},
//		},
//		Sender:           m.Sender,
//		Receiver:         m.Receiver,
//		TimeoutHeight:    clienttypes.Height{
//			RevisionNumber: m.TimeoutHeight.RevisionNumber,
//			RevisionHeight: m.TimeoutHeight.RevisionHeight,
//		},
//		TimeoutTimestamp: m.TimeoutTimestamp,
//	}
//	return ret
//}
// Route implements sdk.Msg
//func (MsgTransfer) Route() string {
//	return RouterKey
//}

// Type implements sdk.Msg
//func (MsgTransfer) Type() string {
//	return TypeMsgTransfer
//}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
//func (msg MsgTransfer) ValidateBasic() error {
//	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
//		return sdkerrors.Wrap(err, "invalid source port ID")
//	}
//	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
//		return sdkerrors.Wrap(err, "invalid source channel ID")
//	}
//	if !msg.Token.IsValid() {
//		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Token.String())
//	}
//	if !msg.Token.IsPositive() {
//		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.Token.String())
//	}
//	// NOTE: sender format must be validated as it is required by the GetSigners function.
//	_, err := sdk.AccAddressFromBech32(msg.Sender)
//	if err != nil {
//		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
//	}
//	if strings.TrimSpace(msg.Receiver) == "" {
//		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
//	}
//	return ValidateIBCDenom(msg.Token.Denom)
//}

// GetSignBytes implements sdk.Msg.

//func (msg MsgTransfer) GetSignBytes() []byte {
//	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
//}

// GetSigners implements sdk.Msg
//func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
//	valAddr, err := sdk.AccAddressFromBech32(msg.Sender)
//	if err != nil {
//		panic(err)
//	}
//	return []sdk.AccAddress{valAddr}
//}
