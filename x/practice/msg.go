package practice

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendCoinsMsg defines a Send Coins Message
type SendCoinsMsg struct {
	FromAddress sdk.AccAddress
	ToAddress   sdk.AccAddress
	Amount      sdk.Coins
}

// NewSendCoinMsg is constructor for SendCoinsMsg
func NewSendCoinMsg(from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) SendCoinsMsg {
	return SendCoinsMsg{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
	}
}

// Route return route name
func (msg SendCoinsMsg) Route() string { return "practice" }

// Type return type
func (msg SendCoinsMsg) Type() string { return "send_coins" }

// ValidateBasic Implements Msg.
func (msg SendCoinsMsg) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() || msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.FromAddress.String())
	}

	if msg.Amount.IsZero() {
		return sdk.ErrInvalidCoins("amount is zero")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg SendCoinsMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg SendCoinsMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
