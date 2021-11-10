package main

import (
	"context"
	"flag"
	"log"
)

func main() {
	var tokenID int64
	flag.Int64Var(&tokenID, "t", -1, "The token ID.")
	flag.Parse()

	if tokenID == -1 {
		log.Fatal("Token ID is required")
	}

	ctx := context.Background()

	tokenURI, err := GetTokenURI(ctx, tokenID)
	if err != nil {
		panic(err)
	}
	log.Printf("Token URI: %s\n", tokenURI)

	ownerAddr, err := GetTokenOwner(ctx, tokenID)
	if err != nil {
		panic(err)
	}
	log.Printf("Owner address: %s\n", ownerAddr)
}
