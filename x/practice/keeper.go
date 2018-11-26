package practice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage
// and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper     bank.Keeper
	amountStoreKey sdk.StoreKey
	cdc            *codec.Codec
}

// NewKeeper creates new instances of keeper
func NewKeeper(coinKeeper bank.Keeper, amountStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:     coinKeeper,
		amountStoreKey: amountStoreKey,
		cdc:            cdc,
	}
}

func (k Keeper) setAmount(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coins) {
	store := ctx.KVStore(k.amountStoreKey)
	store.Set(address, k.cdc.MustMarshalBinaryBare(amount))
}

func (k Keeper) getAmount(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.amountStoreKey)
	bz := store.Get(address)
	var amount sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(bz, &amount)
	return amount
}
