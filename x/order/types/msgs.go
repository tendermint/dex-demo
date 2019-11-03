package types

import (
	"github.com/tendermint/dex-demo/pkg/matcheng"
	"github.com/tendermint/dex-demo/pkg/serde"
	"github.com/tendermint/dex-demo/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgStop struct {
	Post MsgPost
	// The price at the time the tx is signed.
	// Needed to determine which direction the price travels in order to
	// know whether to trigger above or below MsgPost.Price
	InitPrice sdk.Uint
	Relayer   sdk.AccAddress
	RelayFee  sdk.Coins
}

func NewMsgStop(owner sdk.AccAddress, marketID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16, initPrice sdk.Uint, relayer sdk.AccAddress, relayFee sdk.Coins) MsgStop {
	return MsgStop{
		Post: MsgPost{
			Owner:       owner,
			MarketID:    marketID,
			Direction:   direction,
			Price:       price,
			Quantity:    quantity,
			TimeInForce: tif,
		},
		InitPrice: initPrice,
		Relayer:   relayer,
		RelayFee:  relayFee,
	}
}

func (msg MsgStop) Route() string {
	return "order"
}

func (msg MsgStop) Type() string {
	return "stop"
}

func (msg MsgStop) ValidateBasic() sdk.Error {
	// Don't need to check a bool like Buy because t & f should be able to be used here?
	if !msg.Post.MarketID.IsDefined() {
		return sdk.ErrUnauthorized("invalid market ID")
	}
	if msg.Post.Price.IsZero() {
		return sdk.ErrInvalidCoins("price cannot be zero")
	}
	if msg.Post.Quantity.IsZero() {
		return sdk.ErrInvalidCoins("quantity cannot be zero")
	}
	if msg.Post.TimeInForce == 0 {
		return sdk.ErrInternal("time in force cannot be zero")
	}
	if msg.Post.TimeInForce > MaxTimeInForce {
		return sdk.ErrInternal("time in force cannot be larger than 600")
	}
	if msg.InitPrice.IsZero() {
		return sdk.ErrInvalidCoins("pastPrice cannot be zero")
	}
	if msg.RelayFee.IsZero() {
		return sdk.ErrInvalidCoins("relayFee cannot be zero")
	}
	return nil
}

func (msg MsgStop) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgStop) GetSigners() []sdk.AccAddress {
	return msg.Post.GetSigners()
}

type MsgPost struct {
	Owner       sdk.AccAddress
	MarketID    store.EntityID
	Direction   matcheng.Direction
	Price       sdk.Uint
	Quantity    sdk.Uint
	TimeInForce uint16
}

func NewMsgPost(owner sdk.AccAddress, marketID store.EntityID, direction matcheng.Direction, price sdk.Uint, quantity sdk.Uint, tif uint16) MsgPost {
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
	if !msg.MarketID.IsDefined() {
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
	OrderID store.EntityID
}

func NewMsgCancel(owner sdk.AccAddress, orderID store.EntityID) MsgCancel {
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
	if !msg.OrderID.IsDefined() {
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
