package types

import "github.com/cosmos/cosmos-sdk/codec"

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStop{}, "order/Stop", nil)
	cdc.RegisterConcrete(MsgPost{}, "order/Post", nil)
	cdc.RegisterConcrete(MsgCancel{}, "order/Cancel", nil)
}
