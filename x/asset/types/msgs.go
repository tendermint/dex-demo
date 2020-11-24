package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/dex-demo/pkg/serde"
	"github.com/tendermint/dex-demo/types/errs"
)

type MsgCreate struct {
	Name        string
	Symbol      string
	Owner       sdk.AccAddress
	TotalSupply sdk.Uint
}

type MsgMint struct {
	ID     sdk.Uint
	Minter sdk.AccAddress
	Amount sdk.Uint
}

func NewMsgMint(id sdk.Uint, minter sdk.AccAddress, amount sdk.Uint) MsgMint {
	return MsgMint{
		ID:     id,
		Minter: minter,
		Amount: amount,
	}
}

func (msg MsgMint) Route() string {
	return "asset"
}

func (msg MsgMint) Type() string {
	return "mint"
}

func (msg MsgMint) ValidateBasic() sdk.Error {
	if msg.ID.IsZero() {
		return errs.ErrNotFound("asset ID must exist")
	}
	if msg.Minter.Empty() {
		return sdk.ErrInvalidAddress(msg.Minter.String())
	}
	return nil
}

func (msg MsgMint) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Minter}
}

type MsgBurn struct {
	ID     sdk.Uint
	Burner sdk.AccAddress
	Amount sdk.Uint
}

func NewMsgBurn(id sdk.Uint, burner sdk.AccAddress, amount sdk.Uint) MsgBurn {
	return MsgBurn{
		ID:     id,
		Burner: burner,
		Amount: amount,
	}
}

func (msg MsgBurn) Route() string {
	return "asset"
}

func (msg MsgBurn) Type() string {
	return "burn"
}

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if msg.ID.IsZero() {
		return errs.ErrNotFound("asset ID must exist")
	}
	if msg.Burner.Empty() {
		return sdk.ErrInvalidAddress(msg.Burner.String())
	}
	return nil
}

func (msg MsgBurn) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Burner}
}

type MsgTransfer struct {
	ID     sdk.Uint
	From   sdk.AccAddress
	To     sdk.AccAddress
	Amount sdk.Uint
}

func NewMsgTransfer(id sdk.Uint, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Uint) MsgTransfer {
	return MsgTransfer{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
	}
}

func (msg MsgTransfer) Route() string {
	return "asset"
}

func (msg MsgTransfer) Type() string {
	return "transfer"
}

func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if msg.ID.IsZero() {
		return errs.ErrNotFound("asset ID must exist")
	}
	if msg.From.Empty() || msg.To.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return serde.MustMarshalSortedJSON(msg)
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
