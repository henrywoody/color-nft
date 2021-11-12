package provenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/henrywoody/color-nft/ipfs"
	"github.com/henrywoody/color-nft/token"
	"github.com/henrywoody/color-nft/utils"
)

// CalculateProvenanceHashFromTokens calculates the provenance hash for a given
// set of tokens in the given order (should be the original, generated order,
// not the tokenID order).
func CalculateProvenanceHashFromTokens(tokens []*token.Token) string {
	var b strings.Builder
	for _, t := range tokens {
		hash := sha256.Sum256([]byte(t.Image()))
		b.WriteString(hex.EncodeToString(hash[:]))
	}
	hash := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(hash[:])
}

// CalculateProvenanceHashFromIPFSPaths calculates the provenance hash for a
// given set of tokens in the given order (should be in the original, generated
// order, not the tokenID order).
func CalculateProvenanceHashFromIPFSPaths(ctx context.Context, ipfsNode *ipfs.IPFSNode, ipfsPaths []string, outDirPath string) (string, error) {
	var b strings.Builder
	for i, ipfsPath := range ipfsPaths {
		outFilePath := filepath.Join(outDirPath, strconv.Itoa(i))
		if err := ipfsNode.GetFileOrDirectory(ctx, ipfsPath, outFilePath); err != nil {
			return "", fmt.Errorf("error getting image file from IPFS (%s): %v", ipfsPath, err)
		}
		fileData, err := os.ReadFile(outFilePath)
		if err != nil {
			return "", fmt.Errorf("error reading image file: %v", err)
		}
		hash := sha256.Sum256(fileData)
		b.WriteString(hex.EncodeToString(hash[:]))
	}
	hash := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(hash[:]), nil
}

func OriginalIDFromTokenID(tokenID, startIndex int) int {
	return mod(tokenID+startIndex, utils.MaxTokens)
}

func TokenIDFromOriginalID(originalID, startIndex int) int {
	return mod(originalID-startIndex, utils.MaxTokens)
}

func mod(n, m int) int {
	return ((n % m) + m) % m
}
