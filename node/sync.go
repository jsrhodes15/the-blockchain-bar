package node

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"net/http"
	"time"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(45 * time.Second)

	for {
		select {
		case <-ticker.C:
			n.doSync()

		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) doSync() {
	for _, peer := range n.knownPeers {
		if n.ip == peer.IP && n.port == peer.Port {
			continue
		}

		color.Green("Searching for new Peers and their Blocks and Peers: '%s'\n", peer.TcpAddress())

		status, err := queryPeerStatus(peer)
		if err != nil {
			color.Magenta("ERROR: %s\n", err)
			color.Magenta("Peer '%s' was removed from KnownPeers\n", peer.TcpAddress())

			n.RemovePeer(peer)

			continue
		}

		err = n.joinKnownPeers(peer)
		if err != nil {
			color.Magenta("ERROR: %s\n", err)
			continue
		}

		err = n.syncBlocks(peer, status)
		if err != nil {
			color.Magenta("ERROR: %s\n", err)
			continue
		}

		err = n.syncKnownPeers(peer, status)
		if err != nil {
			color.Magenta("ERROR: %s\n", err)
			continue
		}
	}
}

func (n *Node) syncBlocks(peer PeerNode, status StatusRes) error {
	localBlockNumber := n.state.LatestBlock().Header.Number

	// If peer has no blocks, ignore it
	if status.Hash.IsEmpty() {
		return nil
	}

	// If peer has less blocks than us, ignore it
	if status.Number < localBlockNumber {
		return nil
	}

	// If it's the genesis block and we already synced it, ignore it
	if status.Number == 0 && !n.state.LatestBlockHash().IsEmpty() {
		return nil
	}

	// Display found 1 new block if we sync the genesis block (block0)
	newBlocksCount := status.Number - localBlockNumber
	if localBlockNumber == 0 && status.Number == 0 {
		newBlocksCount = 1
	}

	color.Yellow("Found %d new blocks from Peer %s\n", newBlocksCount, peer.TcpAddress())

	blocks, err := fetchBlocksFromPeer(peer, n.state.LatestBlockHash())
	if err != nil {
		return err
	}

	return n.state.AddBlocks(blocks)
}

func (n *Node) syncKnownPeers(peer PeerNode, status StatusRes) error {
	for _, statusPeer := range status.KnownPeers {
		if !n.IsKnownPeer(statusPeer) {
			color.Cyan("Found new Peer %s\n", statusPeer.TcpAddress())

			n.AddPeer(statusPeer)
		}
	}

	return nil
}

func (n *Node) joinKnownPeers(peer PeerNode) error {
	if peer.connected {
		return nil
	}

	url := fmt.Sprintf(
		"http://%s%s?%s=%s&%s=%d",
		peer.TcpAddress(),
		endpointAddPeer,
		endpointAddPeerQueryKeyIP,
		n.ip,
		endpointAddPeerQueryKeyPort,
		n.port,
	)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	addPeerRes := AddPeerRes{}
	err = readRes(res, &addPeerRes)
	if err != nil {
		return err
	}
	if addPeerRes.Error != "" {
		return fmt.Errorf(addPeerRes.Error)
	}

	knownPeer := n.knownPeers[peer.TcpAddress()]
	knownPeer.connected = addPeerRes.Success

	n.AddPeer(knownPeer)

	if !addPeerRes.Success {
		return fmt.Errorf("unable to join KnownPeers of '%s'", peer.TcpAddress())
	}

	return nil
}

func queryPeerStatus(peer PeerNode) (StatusRes, error) {
	url := fmt.Sprintf("http://%s%s", peer.TcpAddress(), endpointStatus)
	res, err := http.Get(url)
	if err != nil {
		return StatusRes{}, err
	}

	statusRes := StatusRes{}
	err = readRes(res, &statusRes)
	if err != nil {
		return StatusRes{}, err
	}

	return statusRes, nil
}

func fetchBlocksFromPeer(peer PeerNode, fromBlock database.Hash) ([]database.Block, error) {
	color.Yellow("Importing blocks from Peer %s...", peer.TcpAddress())

	url := fmt.Sprintf(
		"http://%s%s?%s=%s",
		peer.TcpAddress(),
		endpointSync,
		endpointSyncQueryKeyFromBlock,
		fromBlock.Hex(),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	syncRes := SyncRes{}
	err = readRes(res, &syncRes)
	if err != nil {
		return nil, err
	}

	return syncRes.Blocks, nil
}
