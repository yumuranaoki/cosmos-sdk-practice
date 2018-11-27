package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/tendermint/tendermint/libs/cli"

	app "github.com/yumuranaoki/cosmos-sdk-practice/app"
	practicecmd "github.com/yumuranaoki/cosmos-sdk-practice/x/practice/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

const storeAcc = "acc"

var (
	rootCmd = &cobra.Command{
		Use:   "practicecli",
		Short: "practice client",
	}
	defaultCLIHome = os.ExpandEnv("$HOME/.practicecli")
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(client.ConfigCmd())
	rpc.AddCommands(rootCmd)

	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(queryCmd, cdc)
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		practicecmd.GetCmdBalance("practice", cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		practicecmd.GetCmdSendCoins(cdc),
	)...)

	rootCmd.AddCommand(
		queryCmd,
		txCmd,
		client.LineBreak,
	)

	rootCmd.AddCommand(
		keys.Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "Practice", defaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
