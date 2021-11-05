# Color NFT

A simple ERC721 non-fungible token (NFT).

## ERC721

Standard: https://eips.ethereum.org/EIPS/eip-721

Implementation library: https://github.com/OpenZeppelin/openzeppelin-contracts

## IPFS

Image files and metadata files are stored on [IPFS](http://ipfs.io).

Go repo: https://github.com/ipfs/go-ipfs

Usage as a library example: https://github.com/ipfs/go-ipfs/tree/master/docs/examples/go-ipfs-as-a-library

## Go Etherem

Go Ethereum Book: https://goethereumbook.org/en/

Go Ethereum native bindings: https://geth.ethereum.org/docs/dapp/native-bindings

Console: https://geth.ethereum.org/docs/rpc/server

## Creating an Account

Use the `account` package to create a new Ethereum account, this will print the address and private key for the account.

To set up the account in development mode:

```shell
geth --dev console --http --http.port 3334 --allow-insecure-unlock
```

And follow the `account/dev-setup.js` script.

## Deployment

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