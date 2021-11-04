package token

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

type Token struct {
	name            string
	imagesDirPath   string
	metaDataDirPath string
	metaData        *TokenMetaData
}

type TokenMetaData struct {
	Image string `json:"image"`
	Color string `json:"color"`
}

func NewToken(name, imagesDirPath, metaDataDirPath string) *Token {
	return &Token{
		name:            name,
		imagesDirPath:   imagesDirPath,
		metaDataDirPath: metaDataDirPath,
		metaData: &TokenMetaData{
			Color: getRandomColor(),
		},
	}
}

func (t *Token) Name() string {
	return t.name
}

func (t *Token) MetaData() *TokenMetaData {
	return t.metaData
}

func (t *Token) SetImageURI(imageURI string) {
	t.metaData.Image = imageURI
}

const imageTmpl = `<?xml version="1.0"?>
<svg width="1600" height="900"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<rect x="0" y="0" width="1600" height="900" style="fill: %s;" />
</svg>
`

func (t *Token) GenerateImage() error {
	fileContent := fmt.Sprintf(imageTmpl, t.metaData.Color)
	if err := os.WriteFile(t.ImageFilePath(), []byte(fileContent), 0664); err != nil {
		return fmt.Errorf("failed to create image file: %v", err)
	}
	return nil
}

func (t *Token) WriteMetaData() error {
	metaDataFile, err := os.Create(t.MetaDataFilePath())
	if err != nil {
		return fmt.Errorf("failed to create meta data file: %v", err)
	}

	if err := json.NewEncoder(metaDataFile).Encode(t.metaData); err != nil {
		return fmt.Errorf("failed to write to meta data file: %v", err)
	}

	return nil
}

func (t *Token) ImageFilePath() string {
	return filepath.Join(t.imagesDirPath, t.name+".svg")
}

func (t *Token) MetaDataFilePath() string {
	return filepath.Join(t.metaDataDirPath, t.name+".json")
}

func getRandomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
