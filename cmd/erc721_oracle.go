package main

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	ethOracle "github.com/brennanjl/kwil-extension-tools/evm_oracle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kwilteam/kwil-db/common"
	"github.com/kwilteam/kwil-db/core/types"
)

var (
	//go:embed dev_dao.json
	abiJSON []byte
	// stakingABI is the abi for the Credit event
	transferABI            abi.ABI
	transferEventSignature = "Transfer(address,address,uint256)"
	transferEventTopic     = crypto.Keccak256Hash([]byte(transferEventSignature))
)

func init() {
	var err error
	transferABI, err = abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		panic(err)
	}

	ethOracle.RegisterEthListener(ethOracle.EthListener{
		ContractAddresses: []string{"0x25ed58c027921E14D86380eA2646E3a1B5C55A8b"},
		EventSignatures:   []string{transferEventSignature},
		StartHeight:       13153967, // height the nft was deployed
		ExtensionName:     "erc721_oracle",
		// extensions.erc721_oracle.block_sync_chunk_size
		ConfigName:            "erc721_oracle",
		RequiredConfirmations: 3,
		Resolve: func(ctx context.Context, app *common.App, log ethTypes.Log) error {
			if len(log.Topics) == 0 {
				app.Service.Logger.Error("no event signature")
				return nil
			}

			var data []any
			switch log.Topics[0] {
			case transferEventTopic:
				data, err = transferABI.Unpack("Transfer", log.Data)
				fmt.Println("The data is: ", log.Data)
				if err != nil {
					app.Service.Logger.Error("failed to unpack transfer event", "error", err)
					return nil
				}
			default:
				app.Service.Logger.Error("unknown event signature", "signature", log.Topics[0].Hex())
				return nil
			}

			// go-ethereum decodes uint256 as *big.Int
			// check that the data includes the from, to, and tokenId
			if len(data) != 3 {
				app.Service.Logger.Error("invalid data length", "length", len(data))
				return nil
			}

			// validate that the to address is bytes
			to, ok := data[1].(ethcommon.Address)
			if !ok {
				app.Service.Logger.Error("invalid to address", "to", data[1])
				return nil
			}

			// validate that the tokenId is a *big.Int
			tokenId, ok := data[2].(*big.Int)
			if !ok {
				app.Service.Logger.Error("invalid tokenId", "tokenId", data[2])
				return nil
			}

			return assignOwnership(ctx, app, to.Hex(), tokenId.Int64(), types.NewUUIDV5(log.TxHash[:]))
		},
	})

}
