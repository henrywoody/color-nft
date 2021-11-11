package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
)

const tokenPrice int64 = 50_000_000_000_000_000 // must match that in Token.sol

func Mint(ctx context.Context, numTokens int64, privateKeyHex string) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

	auth, err := c.GetAuth(ctx, privateKeyHex)
	if err != nil {
		return err
	}

	auth.Value = big.NewInt(numTokens * tokenPrice)

	tx, err := instance.Mint(auth, big.NewInt(numTokens))
	if err != nil {
		return fmt.Errorf("error minting: %v", err)
	}

	log.Printf("Mint transaction hash: %s\n", tx.Hash().Hex())

	return nil
}

func OwnerMint(ctx context.Context, numTokens int64) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

	auth, err := c.GetAuth(ctx, os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	tx, err := instance.OwnerMint(auth, big.NewInt(numTokens))
	if err != nil {
		return err
	}

	log.Printf("Withdraw transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
