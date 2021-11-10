package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/henrywoody/color-nft/client"
	"github.com/henrywoody/color-nft/contract"
)

func getContract() (*contract.ColorNFT, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	contractAddrHex := os.Getenv("CONTRACT_ADDRESS")
	instance, err := c.GetContract(contractAddrHex)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func getCallOpts(ctx context.Context) (*bind.CallOpts, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error converting private key hex to ECDSA: %v", err)
	}
	fromAddr, err := client.PrivateKeyToAddress(privateKey)
	if err != nil {
		return nil, err
	}

	callOpts := &bind.CallOpts{From: fromAddr, Context: ctx}
	return callOpts, nil
}
