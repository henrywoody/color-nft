package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/token"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TokenMetaData struct {
	Image string `json:"image"`
	Color string `json:"color"`
}

const imagesDirPath = "images"
const metaDataDirPath = "metadata"

func main() {
	var toAddrHex string
	flag.StringVar(&toAddrHex, "to", "", "The address to send the token to.")
	flag.Parse()

	if toAddrHex == "" {
		log.Fatal("A `to` address is required.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := createTokenData(ctx)
	if err != nil {
		panic(err)
	}
}

func createTokenData(ctx context.Context) (string, error) {
	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to spawn ephemeral node: %s", err)
	}

	mkDirIfNotExists(imagesDirPath)
	mkDirIfNotExists(metaDataDirPath)

	now := time.Now().Unix()
	t := token.NewToken(strconv.Itoa(int(now)), imagesDirPath, metaDataDirPath)

	if err := t.GenerateImage(); err != nil {
		return "", err
	}

	imageURI, err := ipfsNode.AddFile(ctx, t.ImageFilePath())
	if err != nil {
		return "", err
	}
	t.SetImageURI(imageURI)

	if err := t.WriteMetaData(); err != nil {
		return "", err
	}

	metaDataURI, err := ipfsNode.AddFile(ctx, t.MetaDataFilePath())
	if err != nil {
		return "", err
	}

	log.Printf("Created token with metadata URI: %s\n", metaDataURI)

	return metaDataURI, nil
}

func mkDirIfNotExists(dirPath string) error {
	exists, err := checkDirExists(dirPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return os.Mkdir(dirPath, 0777)
}

func checkDirExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
