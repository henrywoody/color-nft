package utils

import "path/filepath"

const ImagesDirPath = "image-files"
const metaDataDirPath = "metadata-files"

var MetaDataDirPath = filepath.Join(metaDataDirPath, "original")
var FinalMetaDataDirPath = filepath.Join(metaDataDirPath, "final")

const MaxTokens = 100 // must match that in Token.sol
