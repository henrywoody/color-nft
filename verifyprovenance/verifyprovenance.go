package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/token"
	"github.com/henrywoody/color-nft/utils"
)

func main() {
	var startIndex int
	var ipfsHash string
	var provenanceHash string
	flag.IntVar(&startIndex, "i", -1, "The start index for remapping.")
	flag.StringVar(&ipfsHash, "ih", "", "The IPFS hash (URI) of the metadata directory.")
	flag.StringVar(&provenanceHash, "ph", "", "The provenance hash to compare with.")
	flag.Parse()

	ctx := context.Background()

	matches, err := VerifyProvenance(ctx, startIndex, ipfsHash, provenanceHash)
	if err != nil {
		panic(err)
	}
	if matches {
		log.Println("Provenance hash is VALID")
	} else {
		log.Println("Provenance hash is INVALID")
	}
}

func VerifyProvenance(ctx context.Context, startIndex int, metaDataDirIPFSHash string, givenProvenance string) (bool, error) {
	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		return false, err
	}

	go func() {
		if err := ipfsNode.ConnectToPeers(ctx); err != nil {
			panic(err)
		}
	}()

	currentTimestamp := time.Now().UnixNano()

	metaDataDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("metadata-%d", currentTimestamp))
	if err := ipfsNode.GetFileOrDirectory(ctx, metaDataDirIPFSHash, metaDataDirPath); err != nil {
		return false, fmt.Errorf("error getting metadata files from IPFS: %v", err)
	}

	dir, err := os.ReadDir(metaDataDirPath)
	if err != nil {
		return false, fmt.Errorf("error reading metadata directory: %v", err)
	}

	metaDataObjects := make([]*token.TokenMetaData, 0, utils.MaxTokens)
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(metaDataDirPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return false, fmt.Errorf("error reading metadata file: %v", err)
		}
		var metaData token.TokenMetaData
		if err := json.NewDecoder(file).Decode(&metaData); err != nil {
			return false, fmt.Errorf("error JSON decoding metadata file: %v", err)
		}
		metaDataObjects = append(metaDataObjects, &metaData)
	}

	imageIPFSPaths := make([]string, utils.MaxTokens)
	for _, metaData := range metaDataObjects {
		tokenID, err := strconv.Atoi(metaData.Name)
		if err != nil {
			return false, fmt.Errorf("error converting metadata name to integer: %v", err)
		}
		originalID := mod(tokenID+startIndex, utils.MaxTokens)
		imageIPFSPaths[originalID] = metaData.Image
	}

	imagesOutDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("images-%d", currentTimestamp))
	utils.MkDirIfNotExists(imagesOutDirPath)

	var concatenatedHashesBuilder strings.Builder
	for i, imageIPFSPath := range imageIPFSPaths {
		outFilePath := filepath.Join(imagesOutDirPath, strconv.Itoa(i))
		if err := ipfsNode.GetFileOrDirectory(ctx, imageIPFSPath, outFilePath); err != nil {
			return false, fmt.Errorf("error getting image file from IPFS (%s): %v", imageIPFSPath, err)
		}
		fileData, err := os.ReadFile(outFilePath)
		if err != nil {
			return false, fmt.Errorf("error reading image file: %v", err)
		}
		hash := sha256.Sum256(fileData)
		concatenatedHashesBuilder.WriteString(hex.EncodeToString(hash[:]))
	}
	concatenatedHash := sha256.Sum256([]byte(concatenatedHashesBuilder.String()))
	provenance := hex.EncodeToString(concatenatedHash[:])

	return provenance == givenProvenance, nil
}

func mod(n, m int) int {
	return ((n % m) + m) % m
}
