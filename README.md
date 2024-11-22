# Kwil ERC721 Oracle Example

This is an example of how to use the Kwil ERC721 Oracle to token gate access based on the ownership of an ERC721 token.

## How it Works

The oracle will listen for the `Transfer` event on the ERC721 contract and will update the access control list based on the ownership of the token. 

The access control list is a schema in `cmd/genesis.kf`. It will be deployed at the network genesis.

A sample schema for checking addresses against ERC721 ownership can be found in `./sample_schema.kf`.

## How to Run

1. In `cmd/erc721_oracle.go`, update the `contractAddress` and `contractStartHeight` to the address of the ERC721 contract you want to listen to and the Ethereum block height at which the contract was deployed.

```go
var (
    contractAddres = "your_contract_address"
    contractStartHeight = 12345678
)
```

2. Build the binary with the oracle.

```bash
go build -o .build/kwild-erc721 ./cmd
```

This will create a binary called `kwil-erc721` in the `.build` directory.

3. Set up the node config files.

```bash
.build/kwild-erc721 admin setup init -o .build/node_config
```

This will create a directory called `node_config` in the `.build` directory.

4. Add an Ethereum RPC in `.build/node_config/config.toml`.

The config must be added to the `[app.extensions]` section. Use a websocket address.

```toml

[app.extensions]

[app.extensions.erc721_oracle]

rpc_url = "wss://mainnet.infura.io/ws/v3/your_api_key"
```

5. Start docker following the standard node [start instructions](https://docs.kwil.com/docs/node/quickstart#start).

6. Start the node.

```bash
.build/kwild-erc721 node -r .build/demo-node
```

Once started, the oracle will sync with the Ethereum network and listen for the `Transfer` event on the ERC721 contract. It will update the access control list based on the ownership of the token.

## How to Test

After the node is running and the oracle is caught up with the Ethereum network, you can do the following to test:

1. Deploy the schema in `./sample_schema.kf`.

You can deploy the schema using the [Kwil CLI](https://docs.kwil.com/docs/ref/kwil-cli/database/deploy).

2. **Verify Write Token Gating**: Call the `add_message` procedure. If the address is not the owner of the ERC721 token, the procedure will fail.

You can call the procedure using the [Kwil CLI](https://docs.kwil.com/docs/ref/kwil-cli/database/execute).

3. **Verify Read Token Gating**: Call the `get_messages` procedure. If the address is not the owner of the ERC721 token, the procedure will fail.

You can call the procedure using the [Kwil CLI](https://docs.kwil.com/docs/ref/kwil-cli/database/call).
