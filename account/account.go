package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := makeAccount(); err != nil {
		panic(err)
	}
}

func makeAccount() error {
	ks := keystore.NewKeyStore("/tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	pw := os.Getenv("ACCOUNT_PASSWORD")
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("error generating key: %v", err)
	}
	account, err := ks.ImportECDSA(privateKey, pw)
	if err != nil {
		return fmt.Errorf("error making account from key: %v", err)
	}

	log.Printf("Account address:\n\t%s\n", account.Address.Hex())

	privateKeyBytes := crypto.FromECDSA(privateKey)
	log.Printf("Private Key:\n\t%s\n", hexutil.Encode(privateKeyBytes)[2:])

	return nil
}
