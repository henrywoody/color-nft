package main

import (
	"context"
	"math/big"
)

func GetTokenURI(ctx context.Context, tokenID int64) (string, error) {
	c, instance, err := getContract()
	if err != nil {
		return "", err
	}
	defer c.Close()
	callOpts, err := getCallOpts(ctx)
	if err != nil {
		return "", err
	}
	tokenURI, err := instance.TokenURI(callOpts, big.NewInt(tokenID))
	if err != nil {
		return "", err
	}
	return tokenURI, nil
}

func GetTokenOwner(ctx context.Context, tokenID int64) (string, error) {
	c, instance, err := getContract()
	if err != nil {
		return "", err
	}
	defer c.Close()
	callOpts, err := getCallOpts(ctx)
	if err != nil {
		return "", err
	}
	addr, err := instance.OwnerOf(callOpts, big.NewInt(tokenID))
	if err != nil {
		return "", err
	}
	return addr.Hex(), nil
}
