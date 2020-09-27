package node

import (
	"context"
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"net/http"
)

const DefaultHttpPort = 8080
const nodeStatusEndpoint = "/node/status"
const nodeSyncEndpoint = "/node/sync"
const endpointSyncQueryKeyFromBlock "fromBlock"

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

func (pn PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

type KnownPeers map[string]PeerNode

type Node struct {
	dataDir    string
	port       uint64
	state      *database.State
	knownPeers KnownPeers
}

func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	// Initialize a new map with only one known peer, the bootstrap node
	knownPeers := make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()] = bootstrap

	return &Node{
		dataDir:    dataDir,
		port:       port,
		knownPeers: knownPeers,
	}
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, isActive bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, isActive}
}

func (n *Node) Run() error {
	ctx := context.Background()
	fmt.Println(fmt.Sprintf("Listening on HTTP port: %d", n.port))

	state, err := database.NewStateFromDisk(n.dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	go n.sync(ctx)

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, state)
	})

	http.HandleFunc("/tx/add", func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, state)
	})

	http.HandleFunc(nodeStatusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})

	http.HandleFunc(nodeSyncEndpoint, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n.dataDir)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)
}