# Color NFT

A simple ERC721 non-fungible token (NFT).

## References

### ERC721

Standard: https://eips.ethereum.org/EIPS/eip-721

Implementation library: https://github.com/OpenZeppelin/openzeppelin-contracts

### IPFS

Image files and metadata files are stored on [IPFS](http://ipfs.io).

Go repo: https://github.com/ipfs/go-ipfs

Usage as a library example: https://github.com/ipfs/go-ipfs/tree/master/docs/examples/go-ipfs-as-a-library

### Go Etherem

Go Ethereum Book: https://goethereumbook.org/en/

Go Ethereum native bindings: https://geth.ethereum.org/docs/dapp/native-bindings

Console: https://geth.ethereum.org/docs/rpc/server

## Guide

### Development

1. Start development Ethereum node.

   ```shell
   geth --dev console --http --http.port 3334 --allow-insecure-unlock
   ```

2. Setup account on the node. See [Managing Accounts](#Managing-Accounts) for details.

3. Compile the contract.

   ```shell
   go run ./compile
   ```

4. Deploy the contract.

   ```shell
   go run ./deploy
   ```

5. Create tokens using the main package.

   ```shell
   go run main.go
   ```

6. Query the contract.

   ```shell
   go build -o interact.exe ./interact # this one takes a while to compile so `go run` is not recommended
   ./interact.exe -t 1
   ```

### Managing Accounts

Use the `account` package to create a new Ethereum account, this will print the address and private key for the account. An existing account may be used instead.

To set up the account in development mode:

```shell
geth --dev console --http --http.port 3334 --allow-insecure-unlock
```

And follow the `account/dev-setup.js` script.

### Deployment

Create an account with `geth`:

```shell
geth account new
```

Fund the account at: https://faucet.ropsten.be or https://faucet.dimensions.network

Run local node:

```shell
geth --ropsten --http --http.port 3334
```

Run deploy script (`deploy/deploy.go`).