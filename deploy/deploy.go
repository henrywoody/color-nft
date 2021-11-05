package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/henrywoody/color-nft/contract"
)

func main() {
	if err := deploy(); err != nil {
		panic(err)
	}
}

var chainIDs = map[string]int64{
	"mainnet":      1,
	"ropsten":      3,
	"rinkeby":      4,
	"goerli":       5,
	"kovan":        42,
	"geth-private": 1337,
}

func deploy() error {
	client, err := ethclient.Dial("http://localhost:3334")
	if err != nil {
		return fmt.Errorf("error dialing client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return fmt.Errorf("error converting private key hex to ECDSA: %v", err)
	}
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}

	ctx := context.Background()
	fromAddr := crypto.PubkeyToAddress(*publicKey)
	log.Printf("Deploying from: %s\n", fromAddr.Hex())

	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("error getting pending nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("error suggesting gas price: %v", err)
	}

	chainID := big.NewInt(chainIDs["geth-private"])
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("error getting transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = uint64(0)  // in units
	auth.GasPrice = gasPrice

	address, tx, _, err := contract.DeployColorNFT(auth, client)
	if err != nil {
		return fmt.Errorf("error deploying contract: %v", err)
	}

	log.Printf("Contract Address: %s\n", address.Hex())
	log.Printf("Transaction Hash: %s\n", tx.Hash().Hex())

	return nil
}
