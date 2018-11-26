package practice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler return a handler constructor
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SendCoinsMsg:
			return handleSendCoinMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleSendCoinMsg handle SendCoinMsg
func handleSendCoinMsg(ctx sdk.Context, keeper Keeper, msg SendCoinsMsg) sdk.Result {
	if msg.FromAddress.Equals(msg.ToAddress) {
		return sdk.ErrUnauthorized("Incorrect sender").Result()
	}

	fromAddressBalance := keeper.getAmount(ctx, msg.FromAddress)
	toAddressBalance := keeper.getAmount(ctx, msg.ToAddress)

	if !fromAddressBalance.IsAllGTE(msg.Amount) {
		return sdk.ErrInsufficientCoins("insufficient").Result()
	}

	newFromAddressBalance := fromAddressBalance.Minus(msg.Amount)
	newToAddressBalance := toAddressBalance.Plus(msg.Amount)

	keeper.setAmount(ctx, msg.FromAddress, newFromAddressBalance)
	keeper.setAmount(ctx, msg.ToAddress, newToAddressBalance)

	return sdk.Result{}
}
