package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	validCommands := []string{
		"`get-provenance-hash`",
		"`get-start-index`",
		"`mint`",
		"`owner-mint`",
		"`set-base-uri`",
		"`set-provenance-hash`",
		"`set-start-index`",
		"`token-info`",
		"`withdraw`",
	}
	validCommandsString := strings.Join(validCommands, ", ")

	var cmd string
	var tokenID int64
	var numTokens int64
	var baseURI string
	var provenanceHash string
	flag.StringVar(&cmd, "c", "", fmt.Sprintf("The command to run. Valid values: %s.", validCommandsString))
	flag.Int64Var(&tokenID, "t", 0, "The token ID (for `token-info`).")
	flag.Int64Var(&numTokens, "n", 0, "The number of tokens to mint (for `mint` and `owner-mint`).")
	flag.StringVar(&baseURI, "u", "", "The new base URI (for `set-base-uri`).")
	flag.StringVar(&provenanceHash, "p", "", "The provenance hash (for `set-provenance-hash`).")
	flag.Parse()

	ctx := context.Background()

	switch cmd {
	case "get-provenance-hash":
		provenanceHash, err := GetProvenanceHash(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("Provenance Hash: %s\n", provenanceHash)
	case "get-start-index":
		startIndex, err := GetStartIndex(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("Start index: %d\n", startIndex)
	case "mint":
		if err := Mint(ctx, numTokens, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "owner-mint":
		if err := OwnerMint(ctx, numTokens); err != nil {
			panic(err)
		}
	case "set-base-uri":
		if err := SetBaseURI(ctx, baseURI, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "set-provenance-hash":
		if err := SetProvenanceHash(ctx, provenanceHash, os.Getenv("PRIVATE_KEY")); err != nil {
			panic(err)
		}
	case "set-start-index":
		if err := SetStartIndex(ctx, os.Getenv("PRIVATE_KEY")); err != nil {
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
		log.Fatalf("Unknown command: '%s'\nValid commands are: %s", cmd, validCommandsString)
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
