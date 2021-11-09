package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/henrywoody/color-nft/contract"
)

type Client struct {
	*ethclient.Client
}

func NewClient() (*Client, error) {
	c, err := ethclient.Dial(os.Getenv("ETH_CLIENT_URL"))
	if err != nil {
		return nil, fmt.Errorf("error dialing client: %v", err)
	}
	return &Client{c}, nil
}

var chainIDs = map[string]int64{
	"mainnet":      1,
	"ropsten":      3,
	"rinkeby":      4,
	"goerli":       5,
	"kovan":        42,
	"geth-private": 1337,
}

func (c *Client) GetAuth(ctx context.Context) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error converting private key hex to ECDSA: %v", err)
	}
	fromAddr, err := PrivateKeyToAddress(privateKey)
	if err != nil {
		return nil, err
	}

	nonce, err := c.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting pending nonce: %v", err)
	}

	gasPrice, err := c.SuggestGasPrice(ctx)
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

func (c *Client) GetContract(contractAddrHex string) (*contract.ColorNFT, error) {
	contractAddr := common.HexToAddress(contractAddrHex)
	instance, err := contract.NewColorNFT(contractAddr, c)
	if err != nil {
		return nil, fmt.Errorf("error getting contract: %v", err)
	}
	return instance, nil
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKey), nil
}
