package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henrywoody/color-nft/client"
	"github.com/henrywoody/color-nft/contract"
)

func main() {
	toAddr := "0xf9b1e692d06c824018a4ba504e9ad7d96821da65"
	tokenURI := "ipfs://QmZpxKWLXHy7Y6PcA3dHwgNe1yqgxi3AQbV4fHyCy56AX1"
	if err := mint(toAddr, tokenURI); err != nil {
		panic(err)
	}
}

func mint(toAddrHex, tokenURI string) error {
	c, err := client.GetClient()
	if err != nil {
		return err
	}
	defer c.Close()

	contractAddrHex := os.Getenv("CONTRACT_ADDRESS")
	contractAddr := common.HexToAddress(contractAddrHex)
	instance, err := contract.NewColorNFT(contractAddr, c)
	if err != nil {
		return fmt.Errorf("error getting contract: %v", err)
	}

	name, err := instance.Name(nil)
	if err != nil {
		return fmt.Errorf("error getting name: %v", err)
	}
	log.Printf("Found contract: %s (%s)\n", name, contractAddrHex)

	ctx := context.Background()
	auth, err := client.GetAuth(ctx, c)

	toAddr := common.HexToAddress(toAddrHex)
	tx, err := instance.Mint(auth, toAddr, tokenURI)
	if err != nil {
		return fmt.Errorf("error minting: %v", err)
	}

	log.Printf("Mint transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
