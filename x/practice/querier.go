package practice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryBalance = "balance"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// pathに何が入るかわからないのでerrorなるかも
	addressString := path[0]
	address, _ := sdk.AccAddressFromHex(addressString)
	amount := keeper.getAmount(ctx, address)
	bz := keeper.cdc.MustMarshalBinaryBare(amount)
	return bz, nil
}
