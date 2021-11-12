package main

import (
	"context"
	"fmt"
	"log"
)

func SetStartIndex(ctx context.Context, privateKeyHex string) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

	auth, err := c.GetAuth(ctx, privateKeyHex)
	if err != nil {
		return err
	}

	tx, err := instance.SetStartIndex(auth)
	if err != nil {
		return fmt.Errorf("error setting start index: %v", err)
	}

	log.Printf("Set start index transaction hash: %s\n", tx.Hash().Hex())

	return nil
}

func GetStartIndex(ctx context.Context) (int64, error) {
	c, instance, err := getContract()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	callOpts, err := getCallOpts(ctx)
	if err != nil {
		return 0, err
	}

	startIndex, err := instance.StartIndex(callOpts)
	if err != nil {
		return 0, err
	}

	return startIndex.Int64(), nil
}
