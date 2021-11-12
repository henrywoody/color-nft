package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/metadatautils"
	"github.com/henrywoody/color-nft/provenance"
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
		log.Println("Provenance Hash is VALID")
	} else {
		log.Println("Provenance Hash is INVALID")
	}
}

func VerifyProvenance(ctx context.Context, startIndex int, metaDataDirIPFSHash string, givenProvenanceHash string) (bool, error) {
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

	metaDataObjects, err := metadatautils.GetMetaDataObjects(metaDataDirPath)
	if err != nil {
		return false, err
	}

	imageIPFSPaths, err := getOrderedImageIPFSPaths(metaDataObjects, startIndex)
	if err != nil {
		return false, err
	}

	imagesOutDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("images-%d", currentTimestamp))
	utils.MkDirIfNotExists(imagesOutDirPath)

	provenanceHash, err := provenance.CalculateProvenanceHashFromIPFSPaths(ctx, ipfsNode, imageIPFSPaths, imagesOutDirPath)
	if err != nil {
		return false, err
	}

	return provenanceHash == givenProvenanceHash, nil
}

func getOrderedImageIPFSPaths(metaDataObjects []*token.TokenMetaData, startIndex int) ([]string, error) {
	imageIPFSPaths := make([]string, utils.MaxTokens)
	for _, metaData := range metaDataObjects {
		tokenID, err := strconv.Atoi(metaData.Name)
		if err != nil {
			return nil, fmt.Errorf("error converting metadata name to integer: %v", err)
		}
		originalID := provenance.OriginalIDFromTokenID(tokenID, startIndex)
		imageIPFSPaths[originalID] = metaData.Image
	}
	return imageIPFSPaths, nil
}
