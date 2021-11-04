package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	now := time.Now().Unix()
	outFilePath := filepath.Join("images", strconv.Itoa(int(now))+".svg")
	generateImage(outFilePath)
}

const imageTmpl = `<?xml version="1.0"?>
<svg width="1600" height="900"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
<rect x="0" y="0" width="1600" height="900" style="fill: %s;" />
</svg>
`

func generateImage(outFilePath string) {
	color := getRandomColor()
	fileContent := fmt.Sprintf(imageTmpl, color)
	os.WriteFile(outFilePath, []byte(fileContent), 0644)
}

func getRandomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}
