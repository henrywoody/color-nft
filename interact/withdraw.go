package main

import (
	"context"
	"log"
	"os"

	"github.com/henrywoody/color-nft/client"
)

func Withdraw(ctx context.Context) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	instance, err := c.GetContract(os.Getenv("CONTRACT_ADDRESS"))
	if err != nil {
		return err
	}

	auth, err := c.GetAuth(ctx, os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	tx, err := instance.Withdraw(auth)
	if err != nil {
		return err
	}

	log.Printf("Withdraw transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
