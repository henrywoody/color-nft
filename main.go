package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

	tokens := make([]*token.Token, maxTokens)
	for i := 0; i < maxTokens; i++ {
		t, err := createToken(ctx, i, ipfsNode)
		if err != nil {
			panic(fmt.Errorf("error creating token %d: %v", i, err))
		}
		tokens[i] = t
	}

	provenance := calculateProvenanceHash(tokens)
	log.Printf("Provenance Hash: %s\n", provenance)

	dirURI, err := ipfsNode.AddDirectory(ctx, metaDataDirPath)
	if err != nil {
		panic(err)
	}

	log.Printf("MetaData IPFS Directory: %s\n", dirURI)
}

func createToken(ctx context.Context, generateID int, ipfsNode *ipfs.IPFSNode) (*token.Token, error) {
	t := token.NewToken(strconv.Itoa(generateID), imagesDirPath, metaDataDirPath)

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
