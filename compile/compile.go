package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/crypto"
)

const srcFilePath = "Token.sol"
const depsDirPath = "sol_deps"

func main() {
	srcPath, err := fetchDependencies()
	if err != nil {
		panic(err)
	}

	contracts, err := compiler.CompileSolidity("solc", srcPath)
	if err != nil {
		panic(err)
	}

	abis := []string{}
	bins := []string{}
	sigs := []map[string]string{}
	types := []string{}
	libs := make(map[string]string)
	aliases := make(map[string]string)

	for name, contract := range contracts {
		abi, err := json.Marshal(contract.Info.AbiDefinition)
		if err != nil {
			panic(err)
		}
		abis = append(abis, string(abi))
		bins = append(bins, contract.Code)
		sigs = append(sigs, contract.Hashes)
		nameParts := strings.Split(name, ":")
		typeName := nameParts[len(nameParts)-1]
		types = append(types, typeName)

		libPattern := crypto.Keccak256Hash([]byte(name)).String()[2:36]
		libs[libPattern] = typeName
	}

	pkgName := "contract"
	code, err := bind.Bind(types, abis, bins, sigs, pkgName, bind.LangGo, libs, aliases)
	if err != nil {
		panic(err)
	}

	outFilePath := filepath.Join("contract", "contract.go")
	if err := os.WriteFile(outFilePath, []byte(code), 0664); err != nil {
		panic(err)
	}
}

func fetchDependencies() (string, error) {
	fileData, err := os.ReadFile(srcFilePath)
	if err != nil {
		return "", err
	}

	importRe := regexp.MustCompile(`import "(.+?)";`)
	lines := strings.Split(string(fileData), "\n")
	importPaths := []string{}
	for _, line := range lines {
		match := importRe.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		importPaths = append(importPaths, match[1])
	}

	gitHubRepoURLRe := regexp.MustCompile(`https://github.com/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+`)
	importRepoURLs := make(map[string]bool)
	for _, importPath := range importPaths {
		match := gitHubRepoURLRe.FindStringSubmatch(importPath)
		if match == nil {
			continue
		}
		importRepoURLs[match[0]] = true
	}

	os.RemoveAll(depsDirPath)
	if err := os.Mkdir(depsDirPath, 0774); err != nil {
		return "", err
	}

	os.Chdir(depsDirPath)
	for importRepoURL := range importRepoURLs {
		if err := exec.Command("git", "clone", importRepoURL).Run(); err != nil {
			return "", err
		}
	}
	os.Chdir("..")

	gitHubOrgURLReString := `https://github.com/[A-Za-z0-9_-]+`
	gitHubOrgURLRe := regexp.MustCompile(gitHubOrgURLReString)
	importGitHubRepoRe := regexp.MustCompile(`import "(` + gitHubOrgURLReString + `)`)
	finalLines := make([]string, len(lines))
	for i, line := range lines {
		if !importGitHubRepoRe.MatchString(line) {
			finalLines[i] = line
			continue
		}
		finalLines[i] = gitHubOrgURLRe.ReplaceAllString(line, "./sol_deps")
	}

	srcFilePathParts := strings.Split(srcFilePath, ".")
	finalSrcFilePath := srcFilePathParts[0] + "-CompileReady." + srcFilePathParts[1]

	finalSrc := strings.Join(finalLines, "\n")
	if err := os.WriteFile(finalSrcFilePath, []byte(finalSrc), 0774); err != nil {
		return "", err
	}

	return finalSrcFilePath, nil
}
