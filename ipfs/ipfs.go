package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/libp2p/go-libp2p-core/peer"
)

type IPFSNode struct {
	api icore.CoreAPI
}

// Spawns a node on the default repo location, if the repo exists
func SpawnNode(ctx context.Context) (*IPFSNode, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, err

	}

	api, err := createNode(ctx, defaultPath)
	if err != nil {
		return nil, err
	}
	return &IPFSNode{api: api}, nil
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func SpawnEphemeralNode(ctx context.Context) (*IPFSNode, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %v", err)
	}

	// Spawning an ephemeral IPFS node
	api, err := createNode(ctx, repoPath)
	if err != nil {
		return nil, err
	}

	return &IPFSNode{api: api}, nil
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %v", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %v", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %v", err)
	}

	return nil
}

func createTempRepo() (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %v", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %v", err)
	}

	return repoPath, nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func (n *IPFSNode) AddFile(ctx context.Context, inputFilePath string) (string, error) {
	fileNode, err := getUnixfsNode(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("Could not get file: %v", err)
	}

	cidFile, err := n.api.Unixfs().Add(ctx, fileNode)
	if err != nil {
		return "", fmt.Errorf("Could not add file: %v", err)
	}

	return cidFile.String(), nil
}

func (n *IPFSNode) AddDirectory(ctx context.Context, inputDirPath string) (string, error) {
	dirNode, err := getUnixfsNode(inputDirPath)
	if err != nil {
		return "", fmt.Errorf("Could not get directory: %v", err)
	}

	cidDir, err := n.api.Unixfs().Add(ctx, dirNode)
	if err != nil {
		return "", fmt.Errorf("Could not add directory: %v", err)
	}

	return cidDir.String(), nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (n *IPFSNode) GetFileOrDirectory(ctx context.Context, ipfsHash string, outPath string) error {
	fileCID := icorepath.New(ipfsHash)

	node, err := n.api.Unixfs().Get(ctx, fileCID)
	if err != nil {
		return err
	}

	return files.WriteTo(node, outPath)
}

var peerNodes = []string{
	// IPFS Bootstrapper nodes.
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

	// IPFS Cluster Pinning nodes
	"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
	"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
	"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
	"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
}

func (n *IPFSNode) ConnectToPeers(ctx context.Context) error {
	peerInfos := make(map[peer.ID]*peer.AddrInfo, len(peerNodes))
	for _, addrStr := range peerNodes {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		storedAddrInfo, ok := peerInfos[addrInfo.ID]
		if !ok {
			storedAddrInfo = &peer.AddrInfo{ID: addrInfo.ID}
			peerInfos[storedAddrInfo.ID] = storedAddrInfo
		}
		storedAddrInfo.Addrs = append(storedAddrInfo.Addrs, addrInfo.Addrs...)
	}

	var wg sync.WaitGroup
	wg.Add(len(peerInfos))
	for _, peerInfo := range peerInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			if err := n.api.Swarm().Connect(ctx, *peerInfo); err != nil {
				log.Printf("failed to connect to %s: %v\n", peerInfo.ID, err)
			}
		}(peerInfo)
	}
	wg.Wait()

	return nil
}
