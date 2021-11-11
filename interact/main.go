package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	validCommands := "`mint`, `owner-mint`, `set-baseuri`, `token-info`, `withdraw`"

	var cmd string
	var tokenID int64
	var numTokens int64
	var baseURI string
	flag.StringVar(&cmd, "c", "", fmt.Sprintf("The command to run. Valid values: %s.", validCommands))
	flag.Int64Var(&tokenID, "t", 0, "The token ID (for `token-info`).")
	flag.Int64Var(&numTokens, "n", 0, "The number of tokens to mint (for `mint` and `owner-mint`).")
	flag.StringVar(&baseURI, "u", "", "The new base URI (for `set-baseuri`).")
	flag.Parse()

	ctx := context.Background()

	switch cmd {
	case "mint":
		if err := Mint(ctx, numTokens, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "owner-mint":
		if err := OwnerMint(ctx, numTokens); err != nil {
			panic(err)
		}
	case "set-baseuri":
		if err := SetBaseURI(ctx, baseURI, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "token-info":
		if err := getTokenInfo(ctx, tokenID); err != nil {
			panic(err)
		}
	case "withdraw":
		if err := Withdraw(ctx); err != nil {
			panic(err)
		}
	default:
		log.Fatalf("Unknown command: '%s'\nValid commands are: %s", cmd, validCommands)
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
