package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	appName = "cosmosSdkPractice"
)

type practiceApp struct {
	*bam.BaseApp
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
	}
	return app
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	return cdc
}
