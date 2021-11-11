package main

import (
	"context"
	"log"
	"os"
)

func Withdraw(ctx context.Context) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

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
