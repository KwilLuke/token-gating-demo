package main

import (
	"context"
	_ "embed"
	"encoding/hex"

	ethOracle "github.com/brennanjl/kwil-extension-tools/evm_oracle"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kwilteam/kwil-db/common"
	"github.com/kwilteam/kwil-db/core/types"
)

var (
	contractAddress     = "0x6f32d1b318877e19bC2d55FE0892F698BA4A4A06" // address of the ERC-721 contract
	contractStartHeight = 7259828  
	// stakingABI is the abi for the Credit event
	transferEventSignature = "Transfer(address,address,uint256)"
	transferEventTopic     = crypto.Keccak256Hash([]byte(transferEventSignature))
)

func init() {
	ethOracle.RegisterEthListener(ethOracle.EthListener{
		ContractAddresses: []string{contractAddress},
		EventSignatures:   []string{transferEventSignature},
		StartHeight:       uint64(contractStartHeight), // height the nft was deployed
		ExtensionName:     "erc721_oracle",
		// extensions.erc721_oracle.block_sync_chunk_size
		ConfigName:            "erc721_oracle",
		RequiredConfirmations: 3,
		Resolve: func(ctx context.Context, app *common.App, log ethTypes.Log, kwilBlock *common.BlockContext) error {
			if len(log.Topics) == 0 {
				app.Service.Logger.Error("no event signature")
				return nil
			}

			if log.Topics[0] != transferEventTopic {
				app.Service.Logger.Error("unknown event signature, expected Transfer", "signature", log.Topics[0].Hex())
				return nil
			}

			// Per the erc721 standard, all 3 topics are indexed (https://eips.ethereum.org/EIPS/eip-721)
			// Therefore, we need to decode the data from the topics, rather than the event log data.
			// See docs here: https://goethereumbook.org/event-read/#topics

			// go-ethereum decodes uint256 as *big.Int
			// check that the data includes the from, to, and tokenId
			if len(log.Topics) != 4 {
				app.Service.Logger.Error("expected Transfer event to have 4 topics", "topics", len(log.Topics))
				return nil
			}

			// the first topic is the event signature
			// The next 3 topics are indexed, so we can decode them directly
			// first indexed is the from address
			// second indexed is the to address
			// third indexed is the tokenId

			// all topics are 32 bytes, but addresses are 20 bytes.
			// it is padded as 0x000000000000000000000000 + address
			toTopic := "0x" + hex.EncodeToString(log.Topics[2][12:])

			return assignOwnership(ctx, app, toTopic, log.Topics[3].Big().Int64(), types.NewUUIDV5(log.TxHash[:]), kwilBlock)
		},
	})

}
