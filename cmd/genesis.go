package main

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"

	"github.com/kwilteam/kwil-db/common"
	"github.com/kwilteam/kwil-db/core/types"
	"github.com/kwilteam/kwil-db/extensions/hooks"
	"github.com/kwilteam/kwil-db/parse"
)

var (
	//go:embed genesis.kf
	genesisSchemaBytes []byte
	genesisSchema      *types.Schema
	// networkCaller ensures that only the oracle can update the dataset
	networkCaller = "token_gating_caller"
)

func init() {
	var err error
	genesisSchema, err = parse.Parse(genesisSchemaBytes)
	if err != nil {
		panic(err)
	}
	genesisSchema.Owner = []byte(networkCaller)

	err = hooks.RegisterGenesisHook("token_gating_schema", func(ctx context.Context, app *common.App, chain *common.ChainContext) error {
		return app.Engine.CreateDataset(&common.TxContext{
			Ctx: ctx,
			BlockContext: &common.BlockContext{
				ChainContext: chain,
				Height:       0,
			},
			Signer: []byte(networkCaller),
			Caller: networkCaller,
			TxID:   "genesis_schema",
		}, app.DB, genesisSchema)
	})
	if err != nil {
		panic(err)
	}
}

// // stake allows the oracle to register seen staking events
func assignOwnership(ctx context.Context, app *common.App, owner string, tokenId int64, resolutionID *types.UUID, kwilBlock *common.BlockContext) error {
	txid := sha256.Sum256(resolutionID.Bytes())
	_, err := app.Engine.Procedure(&common.TxContext{
		Ctx:          ctx,
		BlockContext: kwilBlock,
		Signer:       []byte(networkCaller),
		Caller:       networkCaller,
		TxID:         hex.EncodeToString(txid[:]),
	}, app.DB, &common.ExecutionData{
		Dataset:   genesisSchema.DBID(),
		Procedure: "upsert_owner",
		Args:      []any{owner, tokenId},
	})
	return err
}
