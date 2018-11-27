package cli

import (
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/codec"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/yumuranaoki/cosmos-sdk-practice/x/practice"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCmdSendCoins(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "send-coins [address] [amount]",
		Short: "send coins to address",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext()
			.WithCodec(cdc).WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txBldr := authtxb.NewTxBuilderFormCLI().WIthCodec(cdc)

			if err := cliCtx.EnsureAccountExists(): err != nil {
				return err
			}

			fromAccount, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			toAccount, err := sdk.AccAddressFromHex(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := practice.NewSendCoinMsg(fromAccount, toAccount ,coins)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}