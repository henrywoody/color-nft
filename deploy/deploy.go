package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/henrywoody/color-nft/client"
	"github.com/henrywoody/color-nft/contract"
)

func main() {
	if err := deploy(); err != nil {
		panic(err)
	}
}

func deploy() error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.Background()
	auth, err := c.GetAuth(ctx, os.Getenv("PRIVATE_KEY"))

	address, tx, _, err := contract.DeployColorNFT(auth, c)
	if err != nil {
		return fmt.Errorf("error deploying contract: %v", err)
	}

	log.Printf("Contract Address: %s\n", address.Hex())
	log.Printf("Transaction Hash: %s\n", tx.Hash().Hex())

	return nil
}
