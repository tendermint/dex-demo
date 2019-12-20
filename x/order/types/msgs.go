package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/pkg/serde"
)

type MsgPost struct {
	Owner       sdk.AccAddress
	MarketID    sdk.Uint
	Direction   matcheng.Direction
	Price       sdk.Uint
	Quantity    sdk.Uint
	TimeInForce uint16
}

func NewMsgPost(owner sdk.AccAddress, marketID sdk.Uint, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16) MsgPost {
	return MsgPost{
		Owner:       owner,
		MarketID:    marketID,
		Direction:   direction,
		Price:       price,
		Quantity:    quantity,
		TimeInForce: tif,
	}
}

func (msg MsgPost) Route() string {
	return "order"
}

func (msg MsgPost) Type() string {
	return "post"
}

func (msg MsgPost) ValidateBasic() sdk.Error {
	if msg.MarketID.IsZero() {
		return sdk.ErrUnauthorized("invalid market ID")
	}
	if msg.Price.IsZero() {
		return sdk.ErrInvalidCoins("price cannot be zero")
	}
	if msg.Quantity.IsZero() {
		return sdk.ErrInvalidCoins("quantity cannot be zero")
	}
	if msg.TimeInForce == 0 {
		return sdk.ErrInternal("time in force cannot be zero")
	}
	if msg.TimeInForce > MaxTimeInForce {
		return sdk.ErrInternal("time in force cannot be larger than 600")
	}
	return nil
}

func (msg MsgPost) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgPost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgCancel struct {
	Owner   sdk.AccAddress
	OrderID sdk.Uint
}

func NewMsgCancel(owner sdk.AccAddress, orderID sdk.Uint) MsgCancel {
	return MsgCancel{
		Owner:   owner,
		OrderID: orderID,
	}
}

func (msg MsgCancel) Route() string {
	return "order"
}

func (msg MsgCancel) Type() string {
	return "cancel"
}

func (msg MsgCancel) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrUnauthorized("owner cannot be empty")
	}
	if msg.OrderID.IsZero() {
		return sdk.ErrInternal("invalid order ID")
	}
	return nil
}

func (msg MsgCancel) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgCancel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
