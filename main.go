package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	}

	mkDirIfNotExists(imagesDirPath)
	mkDirIfNotExists(metaDataDirPath)

	now := time.Now().Unix()
	t := token.NewToken(strconv.Itoa(int(now)), imagesDirPath, metaDataDirPath)

	if err := t.GenerateImage(); err != nil {
		panic(err)
	}

	imageURI, err := ipfsNode.AddFile(ctx, t.ImageFilePath())
	if err != nil {
		panic(err)
	}
	t.SetImageURI(imageURI)

	if err := t.WriteMetaData(); err != nil {
		panic(err)
	}

	metaDataURI, err := ipfsNode.AddFile(ctx, t.MetaDataFilePath())
	if err != nil {
		panic(err)
	}

	log.Printf("Metadata URI: %s\n", metaDataURI)
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
