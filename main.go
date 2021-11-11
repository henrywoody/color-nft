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

const imagesDirPath = "images"
const metaDataDirPath = "metadata"
const maxTokens = 100 // must match that in Token.sol

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mkDirIfNotExists(imagesDirPath)
	mkDirIfNotExists(metaDataDirPath)

	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		panic(fmt.Errorf("error spawning IPFS node: %v", err))
	}

	for i := 0; i < maxTokens; i++ {
		if err := createTokenData(ctx, i, ipfsNode); err != nil {
			panic(fmt.Errorf("error creating token %d: %v", i, err))
		}
	}

	dirURI, err := ipfsNode.AddDirectory(ctx, metaDataDirPath)
	if err != nil {
		panic(err)
	}

	log.Printf("MetaData IPFS Directory: %s\n", dirURI)
}

func createTokenData(ctx context.Context, generateID int, ipfsNode *ipfs.IPFSNode) error {
	t := token.NewToken(strconv.Itoa(generateID), imagesDirPath, metaDataDirPath)

	if err := t.GenerateImage(); err != nil {
		return err
	}

	imageURI, err := ipfsNode.AddFile(ctx, t.ImageFilePath())
	if err != nil {
		return err
	}
	t.SetImageURI(imageURI)

	if err := t.WriteMetaData(); err != nil {
		return err
	}

	return nil
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
