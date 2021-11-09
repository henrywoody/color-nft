package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial("http://localhost:3334")
	if err != nil {
		return nil, fmt.Errorf("error dialing client: %v", err)
	}
	return client, nil
}

var chainIDs = map[string]int64{
	"mainnet":      1,
	"ropsten":      3,
	"rinkeby":      4,
	"goerli":       5,
	"kovan":        42,
	"geth-private": 1337,
}

func GetAuth(ctx context.Context, client *ethclient.Client) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error converting private key hex to ECDSA: %v", err)
	}
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddr := crypto.PubkeyToAddress(*publicKey)

	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting pending nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error suggesting gas price: %v", err)
	}

	chainID := big.NewInt(chainIDs["geth-private"])
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("error getting transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = uint64(0)  // in units
	auth.GasPrice = gasPrice

	return auth, nil
}
