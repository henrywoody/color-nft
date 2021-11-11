package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/token"
	"github.com/henrywoody/color-nft/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := generateTokens(ctx); err != nil {
		panic(err)
	}
}

func generateTokens(ctx context.Context) error {
	utils.RecursiveMkDirIfNotExists(utils.ImagesDirPath)
	utils.RecursiveMkDirIfNotExists(utils.MetaDataDirPath)

	ipfsNode, err := ipfs.SpawnNode(ctx)
	if err != nil {
		return fmt.Errorf("error spawning IPFS node: %v", err)
	}

	tokens := make([]*token.Token, utils.MaxTokens)
	for i := 0; i < utils.MaxTokens; i++ {
		t, err := createToken(ctx, i, ipfsNode)
		if err != nil {
			return fmt.Errorf("error creating token %d: %v", i, err)
		}
		tokens[i] = t
	}

	provenance := calculateProvenanceHash(tokens)
	log.Printf("Provenance Hash: %s\n", provenance)

	return nil
}

func createToken(ctx context.Context, generateID int, ipfsNode *ipfs.IPFSNode) (*token.Token, error) {
	t := token.NewToken(strconv.Itoa(generateID), utils.ImagesDirPath, utils.MetaDataDirPath)

	if err := t.GenerateImage(); err != nil {
		return nil, err
	}

	imageURI, err := ipfsNode.AddFile(ctx, t.ImageFilePath())
	if err != nil {
		return nil, err
	}
	t.SetImageURI(imageURI)

	if err := t.WriteMetaData(); err != nil {
		return nil, err
	}

	return t, nil
}

func calculateProvenanceHash(tokens []*token.Token) string {
	var b strings.Builder
	for _, t := range tokens {
		hash := sha256.Sum256([]byte(t.Image()))
		b.WriteString(hex.EncodeToString(hash[:]))
	}
	hash := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(hash[:])
}
