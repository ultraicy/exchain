package keeper

import (
	"context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/application/transfer/types"
)

//var _ types.MsgServer = Keeper{}

// See createOutgoingPacket in spec:https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay

type MsgTransferResponse struct {
}

// Transfer defines a rpc handler method for MsgTransfer.
func (k Keeper) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if err := k.SendTransfer(
		ctx, msg.SourcePort, msg.SourceChannel, sdk.Coin{
			Denom: msg.Token.Denom,
			Amount: msg.Token.Amount.ToDec(),
		}, sender, msg.Receiver, msg.TimeoutHeight, msg.TimeoutTimestamp,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC fungible token transfer", "token", msg.Token.Denom, "amount", msg.Token.Amount.String(), "sender", msg.Sender, "receiver", msg.Receiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &MsgTransferResponse{}, nil
}
