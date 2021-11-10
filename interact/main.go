package main

import (
	"context"
	"flag"
	"log"
	"os"
)

func main() {
	var cmd string
	var tokenID int64
	var numTokens int64
	flag.StringVar(&cmd, "c", "", "The command to run. Valid values: `token-info`, `mint`, `owner-mint`, `withdraw`.")
	flag.Int64Var(&tokenID, "t", 0, "The token ID (for `token-info`).")
	flag.Int64Var(&numTokens, "n", 0, "The number of tokens to mint (for `mint` and `owner-mint`).")
	flag.Parse()

	ctx := context.Background()

	switch cmd {
	case "token-info":
		if err := getTokenInfo(ctx, tokenID); err != nil {
			panic(err)
		}
	case "mint":
		if err := Mint(ctx, numTokens, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "owner-mint":
		if err := OwnerMint(ctx, numTokens); err != nil {
			panic(err)
		}
	case "withdraw":
		if err := Withdraw(ctx); err != nil {
			panic(err)
		}
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}

func getTokenInfo(ctx context.Context, tokenID int64) error {
	log.Printf("Token ID: %d\n", tokenID)

	tokenURI, err := GetTokenURI(ctx, tokenID)
	if err != nil {
		return err
	}
	log.Printf("Token URI: %s\n", tokenURI)

	ownerAddr, err := GetTokenOwner(ctx, tokenID)
	if err != nil {
		return err
	}
	log.Printf("Owner address: %s\n", ownerAddr)

	return nil
}
