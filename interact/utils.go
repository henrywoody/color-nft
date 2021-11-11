package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/henrywoody/color-nft/client"
	"github.com/henrywoody/color-nft/contract"
)

func getContract() (*client.Client, *contract.ColorNFT, error) {
	c, err := client.NewClient()
	if err != nil {
		return nil, nil, err
	}

	contractAddrHex := os.Getenv("CONTRACT_ADDRESS")
	instance, err := c.GetContract(contractAddrHex)
	if err != nil {
		return nil, nil, err
	}

	name, err := instance.Name(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting name: %v", err)
	}
	log.Printf("Found contract: %s (%s)\n", name, contractAddrHex)

	return c, instance, nil
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
