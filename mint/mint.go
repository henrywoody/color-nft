package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/henrywoody/color-nft/client"
)

const tokenPrice int64 = 50_000_000_000_000_000 // must match that in Token.sol

func main() {
	ctx := context.Background()
	numTokens := 1
	if err := mint(ctx, numTokens, os.Getenv("PRIVATE_KEY")); err != nil {
		panic(err)
	}
}

func mint(ctx context.Context, numTokens int, privateKeyHex string) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	contractAddrHex := os.Getenv("CONTRACT_ADDRESS")
	instance, err := c.GetContract(contractAddrHex)
	if err != nil {
		return err
	}

	name, err := instance.Name(nil)
	if err != nil {
		return fmt.Errorf("error getting name: %v", err)
	}
	log.Printf("Found contract: %s (%s)\n", name, contractAddrHex)

	auth, err := c.GetAuth(ctx, privateKeyHex)
	if err != nil {
		return err
	}

	auth.Value = big.NewInt(int64(numTokens) * tokenPrice)

	tx, err := instance.Mint(auth, big.NewInt(int64(numTokens)))
	if err != nil {
		return fmt.Errorf("error minting: %v", err)
	}

	log.Printf("Mint transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
