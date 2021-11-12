package metadatautils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/henrywoody/color-nft/token"
	"github.com/henrywoody/color-nft/utils"
)

func GetMetaDataObjects(metaDataDirPath string) ([]*token.TokenMetaData, error) {
	dir, err := os.ReadDir(metaDataDirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata directory: %v", err)
	}

	metaDataObjects := make([]*token.TokenMetaData, 0, utils.MaxTokens)
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(metaDataDirPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error reading metadata file: %v", err)
		}
		var metaData token.TokenMetaData
		if err := json.NewDecoder(file).Decode(&metaData); err != nil {
			return nil, fmt.Errorf("error JSON decoding metadata file: %v", err)
		}
		metaDataObjects = append(metaDataObjects, &metaData)
	}
	return metaDataObjects, nil
}
