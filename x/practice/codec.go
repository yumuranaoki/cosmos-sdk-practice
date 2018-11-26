package practice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(SendCoinsMsg, "practice/SendCoin", nil)
}
