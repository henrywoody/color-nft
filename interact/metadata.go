package main

import (
	"context"
	"log"
)

func SetBaseURI(ctx context.Context, newBaseURI string, privateKeyHex string) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

	auth, err := c.GetAuth(ctx, privateKeyHex)
	if err != nil {
		return err
	}

	tx, err := instance.SetBaseURI(auth, newBaseURI)
	if err != nil {
		return err
	}

	log.Printf("BaseURI update transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
