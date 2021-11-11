package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/token"
	"github.com/henrywoody/color-nft/utils"
)

func main() {
	var startIndex int
	flag.IntVar(&startIndex, "i", -1, "The start index for remapping.")
	flag.Parse()

	if err := remapMetaDataFiles(startIndex); err != nil {
		panic(err)
	}

	ctx := context.Background()
	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		panic(fmt.Errorf("error spawning IPFS node: %v", err))
	}

	dirURI, err := ipfsNode.AddDirectory(ctx, utils.FinalMetaDataDirPath)
	if err != nil {
		panic(fmt.Errorf("error adding metadata directory to IPFS: %v", err))
	}
	log.Printf("MetaData IPFS Directory: %s\n", dirURI)
}

func remapMetaDataFiles(startIndex int) error {
	if !(startIndex >= 0 && startIndex < utils.MaxTokens) {
		return fmt.Errorf("invalid startIndex, must be on interval [0, %d)", utils.MaxTokens)
	}

	utils.RecursiveMkDirIfNotExists(utils.FinalMetaDataDirPath)

	dir, err := os.ReadDir(utils.MetaDataDirPath)
	if err != nil {
		return err
	}

	metaDataObjects := make([]*token.TokenMetaData, 0, utils.MaxTokens)

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(utils.MetaDataDirPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		var metaData token.TokenMetaData
		if err := json.NewDecoder(file).Decode(&metaData); err != nil {
			return err
		}
		metaDataObjects = append(metaDataObjects, &metaData)
	}

	for _, metaData := range metaDataObjects {
		originalID, err := strconv.Atoi(metaData.Name)
		if err != nil {
			return err
		}
		tokenID := mod(originalID-startIndex, utils.MaxTokens)
		metaData.Name = strconv.Itoa(tokenID)

		file, err := os.Create(filepath.Join(utils.FinalMetaDataDirPath, metaData.Name))
		if err != nil {
			return err
		}
		if err := json.NewEncoder(file).Encode(metaData); err != nil {
			return err
		}
	}

	return nil
}

func mod(n, m int) int {
	return ((n % m) + m) % m
}
