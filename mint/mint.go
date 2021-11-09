package mint

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henrywoody/color-nft/client"
)

func Mint(ctx context.Context, toAddrHex, tokenURI string) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	contractAddrHex := os.Getenv("CONTRACT_ADDRESS")
	instance, err := c.GetContract(contractAddrHex)
	if err != nil {
		return err
	}

	name, err := instance.Name(nil)
	if err != nil {
		return fmt.Errorf("error getting name: %v", err)
	}
	log.Printf("Found contract: %s (%s)\n", name, contractAddrHex)

	auth, err := c.GetAuth(ctx)

	toAddr := common.HexToAddress(toAddrHex)
	tx, err := instance.Mint(auth, toAddr, tokenURI)
	if err != nil {
		return fmt.Errorf("error minting: %v", err)
	}

	log.Printf("Minted token to address: %s\n", toAddrHex)
	log.Printf("Mint transaction hash: %s\n", tx.Hash().Hex())

	return nil
}
