package types

import "github.com/cosmos/cosmos-sdk/codec"

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreate{}, "asset/Create", nil)
	cdc.RegisterConcrete(MsgMint{}, "asset/Mint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "asset/Burn", nil)
	cdc.RegisterConcrete(MsgTransfer{}, "asset/Transfer", nil)
}
