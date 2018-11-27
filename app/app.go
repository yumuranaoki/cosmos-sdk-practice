package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/yumuranaoki/cosmos-sdk-practice/x/practice"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	appName = "practice"
)

type practiceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain           *sdk.KVStoreKey
	keyAccount        *sdk.KVStoreKey
	keyPracticeAmount *sdk.KVStoreKey
	keyFeeCollection  *sdk.KVStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	practiceKeeper      practice.Keeper
}

// NewPracticeApp is constructor function for practiceApp
func NewPracticeApp(logger log.Logger, db dbm.DB) *practiceApp {
	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))
	var app = &practiceApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:           sdk.NewKVStoreKey("main"),
		keyAccount:        sdk.NewKVStoreKey("acc"),
		keyPracticeAmount: sdk.NewKVStoreKey("practice_amount"),
		keyFeeCollection:  sdk.NewKVStoreKey("fee_collection"),
	}

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper)

	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	app.practiceKeeper = practice.NewKeeper(
		app.bankKeeper,
		app.keyPracticeAmount,
		app.cdc,
	)

	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	app.Router()
	.AddRoute("bank", bank.NewHandler(app.bankKeeper))
	.AddRoute("practice", practice.NewHandler(app.practiceKeeper))

	app.QueryRouter()
	.AddRoute("practice", practice.NewQuerier(app.practiceKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyPracticeAmount,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

type GenesisState struct {
	Accounts []auth.BaseAccount `json:"accounts"`
}

func (app *practiceApp) initChainer(ctx sdk.Contextm req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, &acc)
	}

	return abci.ResponseInitChain{}
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	practice.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
