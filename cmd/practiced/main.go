package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/common"
	app "github.com/yumuranaoki/cosmos-sdk-practice/app"
)

// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
var DefaultNodeHome = os.ExpandEnv("$HOME/.nameserviced")

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	appInit := server.AppInit{
		AppGenState: server.SimpleAppGenState,
	}

	rootCmd := &cobra.Command{
		Use:               "practiced",
		Short:             "practice app daemon",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(InitCmd(ctx, cdc, appInit))
	server.AddCommands(ctx, cdc, rootCmd, appInit, newApp, exportAppStateAndTMValidators)
	executor := cli.PrepareBaseCmd(rootCmd, "Practice", DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

// newApp is constructor for app
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewPracticeApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	return nil, nil, nil
}

func InitCmd(ctx *server.Context, cdc *codec.Codec, appInit server.AppInit) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
			if err != nil {
				return err
			}
			nodeID := string(nodeKey.ID())

			pk := gaiaInit.ReadOrCreatePrivValidator(config.PrivValidatorFile())
			genTx, appMessage, validator, err := server.SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := appInit.AppGenState(cdc, []json.RawMessage{genTx})
			if err != nil {
				return err
			}
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return gaiaInit.WriteGenesisFile(config.GenesisFile(), chainID, []tmtypes.GenesisValidator{validator}, appStateJSON)
		},
	}
}
