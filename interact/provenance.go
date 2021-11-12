package main

import (
	"context"
	"fmt"
	"log"
)

func SetProvenanceHash(ctx context.Context, provenanceHash, privateKeyHex string) error {
	c, instance, err := getContract()
	if err != nil {
		return err
	}
	defer c.Close()

	auth, err := c.GetAuth(ctx, privateKeyHex)
	if err != nil {
		return err
	}

	tx, err := instance.SetProvenanceHash(auth, provenanceHash)
	if err != nil {
		return fmt.Errorf("error setting start index: %v", err)
	}

	log.Printf("Set provenance hash transaction hash: %s\n", tx.Hash().Hex())

	return nil
}

func GetProvenanceHash(ctx context.Context) (string, error) {
	c, instance, err := getContract()
	if err != nil {
		return "", err
	}
	defer c.Close()

	callOpts, err := getCallOpts(ctx)
	if err != nil {
		return "", err
	}

	provenanceHash, err := instance.ProvenanceHash(callOpts)
	if err != nil {
		return "", err
	}

	return provenanceHash, nil
}
