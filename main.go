package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/henrywoody/color-nft/ipfs"
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
	imageFilePath := filepath.Join(imagesDirPath, fmt.Sprintf("%d.svg", now))
	metaDataFilePath := filepath.Join(metaDataDirPath, fmt.Sprintf("%d.json", now))

	color, err := generateImage(imageFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to create image: %v", err))
	}

	imageURI, err := ipfsNode.AddFile(ctx, imageFilePath)
	if err != nil {
		panic(err)
	}

	metaData := TokenMetaData{
		Image: imageURI,
		Color: color,
	}

	metaDataFile, err := os.Create(metaDataFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to create meta data file: %v", err))
	}

	if err := json.NewEncoder(metaDataFile).Encode(metaData); err != nil {
		panic(fmt.Errorf("failed to write to meta data file: %v", err))
	}

	metaDataURI, err := ipfsNode.AddFile(ctx, metaDataFilePath)
	if err != nil {
		panic(err)
	}

	log.Printf("Metadata URI: %s\n", metaDataURI)
}

const imageTmpl = `<?xml version="1.0"?>
<svg width="1600" height="900"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<rect x="0" y="0" width="1600" height="900" style="fill: %s;" />
</svg>
`

func generateImage(outFilePath string) (string, error) {
	color := getRandomColor()
	fileContent := fmt.Sprintf(imageTmpl, color)
	err := os.WriteFile(outFilePath, []byte(fileContent), 0664)
	if err != nil {
		return "", err
	}
	return color, nil
}

func getRandomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
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
