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

Note that the `interact` script takes a bit of time to build and might be used frequently, it is therefore recommended to create a build of the `interact` package before use (rather than using `go run`).

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

5. Create tokens using the generate package.

   ```shell
   go run ./generate
   # => Provenance Hash: 9882878611ce91e786840e69c91bef706013e0a4bb49a46685b954228d2cd8d2
   ```

6. Set the provenance on the contract using the output of the previous step.

   ```shell
   ./interact.exe -c set-provenance-hash -p 9882878611ce91e786840e69c91bef706013e0a4bb49a46685b954228d2cd8d2
   ```

7. Mint tokens.

   ```shell
   ./interact.exe -c mint -n 20
   ```

8. Once the `startIndexBlockNumber` has been set on the contract (first mint after the `revealTimestamp` or on the last mint when max tokens has been reached), set the `startIndex` on the contract.

   ```shell
   ./interact.exe -c set-start-index
   ```

9. Get the `startIndex` from the contract and use that to remap the token metadata files add the metadata files to IPFS.

   ```shell
   ./interact.exe -c get-start-index
   # => Start index: 42
   go run ./metadata -i 42
   # => MetaData IPFS Directory: /ipfs/QmYzP8bKDv94XmrVAacDzG1w5fHjuMh9uzCjbgn8HRzepT
   ```

10. Set the `baseURI` on the contract using the output from the previous step (make sure to add a trailing slash to the URI).

    ```shell
    ./interact.exe -c set-base-uri -u "/ipfs/QmYzP8bKDv94XmrVAacDzG1w5fHjuMh9uzCjbgn8HRzepT/"
    ```

11. Query the contract for info on a token.

    ```shell
    ./interact.exe -c token-info -t 0
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