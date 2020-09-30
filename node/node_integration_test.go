package node

import (
	"context"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"github.com/jsrhodes15/the-blockchain-bar/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNode_Run(t *testing.T) {
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	n := New(datadir, "127.0.0.1", 8085, database.NewAccount("jrhodes"), PeerNode{})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err = n.Run(ctx)
	if err != nil {
		t.Fatal("node server was suppose to close after 5s")
	}
}

func TestNode_Mining(t *testing.T) {
	// Remove the test directory if it already exists
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)

	// Construct a new Node instance and configure
	// Jrhodes as a miner
	n := New(datadir, nInfo.IP, nInfo.Port, database.NewAccount("jrhodes"), nInfo)

	// Allow the mining to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(
		context.Background(),
		time.Minute*30,
	)

	// Schedule a new TX in 3 seconds from now, in a separate thread
	// because the n.Run() few lines below is a blocking call
	go func() {
		time.Sleep(time.Second * miningIntervalSeconds / 3)
		tx := database.NewTx("jrhodes", "meads", 1, "")

		_ = n.AddPendingTX(tx, nInfo)
	}()

	// Schedule a new TX in 12 seconds from now simulating
	// that it came in - while the first TX is being mined
	go func() {
		time.Sleep(time.Second*miningIntervalSeconds + 2)
		tx := database.NewTx("jrhodes", "meads", 2, "")

		_ = n.AddPendingTX(tx, nInfo)
	}()

	go func() {
		// Periodically check if we mined the 2 blocks
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	// Run the node, mining and everything in a blocking call (hence the go-routines before)
	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("2 pending TX not mined into 2 under 30m")
	}
}

func TestNode_MiningStopsOnNewSyncedBlock(t *testing.T) {
	// Remove the test directory if it already exists
	datadir := getTestDataDirPath()
	err := fs.RemoveDir(datadir)
	if err != nil {
		t.Fatal(err)
	}

	// Required for AddPendingTX() to describe
	// from what node the TX came from (local node in this case)
	nInfo := NewPeerNode(
		"127.0.0.1",
		8085,
		false,
		database.NewAccount(""),
		true,
	)

	jrhodesAcc := database.NewAccount("jrhodes")
	meadsAcc := database.NewAccount("meads")

	n := New(datadir, nInfo.IP, nInfo.Port, meadsAcc, nInfo)

	// Allow the test to run for 30 mins, in the worst case
	ctx, closeNode := context.WithTimeout(context.Background(), time.Minute*30)

	tx1 := database.NewTx("jrhodes", "meads", 1, "")
	tx2 := database.NewTx("jrhodes", "meads", 2, "")
	tx2Hash, _ := tx2.Hash()

	// Pre-mine a valid block without running the `n.Run()`
	// with Jrhodes as a miner who will receive the block reward,
	// to simulate the block came on the fly from another peer
	validPreMinedPb := NewPendingBlock(database.Hash{}, 0, jrhodesAcc, []database.Tx{tx1})
	validSyncedBlock, err := Mine(ctx, validPreMinedPb)
	if err != nil {
		t.Fatal(err)
	}

	// Add 2 new TXs into the Meads's node
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds - 2))

		err := n.AddPendingTX(tx1, nInfo)
		if err != nil {
			t.Fatal(err)
		}

		err = n.AddPendingTX(tx2, nInfo)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Once the Meads is mining the block, simulate that
	// Jrhodes mined the block with TX1 in it faster
	go func() {
		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining")
		}

		_, err := n.state.AddBlock(validSyncedBlock)
		if err != nil {
			t.Fatal(err)
		}
		// Mock the Jrhodes's block came from a network
		n.newSyncedBlocks <- validSyncedBlock

		time.Sleep(time.Second * 2)
		if n.isMining {
			t.Fatal("synced block should have canceled mining")
		}

		// Mined TX1 by Jrhodes should be removed from the Mempool
		_, onlyTX2IsPending := n.pendingTXs[tx2Hash.Hex()]

		if len(n.pendingTXs) != 1 && !onlyTX2IsPending {
			t.Fatal("synced block should have canceled mining of already mined TX")
		}

		time.Sleep(time.Second * (miningIntervalSeconds + 2))
		if !n.isMining {
			t.Fatal("should be mining again the 1 TX not included in synced block")
		}
	}()

	go func() {
		// Regularly check whenever both TXs are now mined
		ticker := time.NewTicker(time.Second * 10)

		for {
			select {
			case <-ticker.C:
				if n.state.LatestBlock().Header.Number == 1 {
					closeNode()
					return
				}
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)

		// Take a snapshot of the DB balances
		// before the mining is finished and the 2 blocks
		// are created.
		startingJrhodesBalance := n.state.Balances[jrhodesAcc]
		startingMeadsBalance := n.state.Balances[meadsAcc]

		// Wait until the 30 mins timeout is reached or
		// the 2 blocks got already mined and the closeNode() was triggered
		<-ctx.Done()

		endJrhodesBalance := n.state.Balances[jrhodesAcc]
		endMeadsBalance := n.state.Balances[meadsAcc]

		// In TX1 Jrhodes transferred 1 TBB token to Meads
		// In TX2 Jrhodes transferred 2 TBB tokens to Meads
		expectedEndJrhodesBalance := startingJrhodesBalance - tx1.Value - tx2.Value + database.BlockReward
		expectedEndMeadsBalance := startingMeadsBalance + tx1.Value + tx2.Value + database.BlockReward

		if endJrhodesBalance != expectedEndJrhodesBalance {
			t.Fatalf("Jrhodes expected end balance is %d not %d", expectedEndJrhodesBalance, endJrhodesBalance)
		}

		if endMeadsBalance != expectedEndMeadsBalance {
			t.Fatalf("Meads expected end balance is %d not %d", expectedEndMeadsBalance, endMeadsBalance)
		}

		t.Logf("Starting Jrhodes balance: %d", startingJrhodesBalance)
		t.Logf("Starting Meads balance: %d", startingMeadsBalance)
		t.Logf("Ending Jrhodes balance: %d", endJrhodesBalance)
		t.Logf("Ending Meads balance: %d", endMeadsBalance)
	}()

	_ = n.Run(ctx)

	if n.state.LatestBlock().Header.Number != 1 {
		t.Fatal("was suppose to mine 2 pending TX into 2 valid blocks under 30m")
	}

	if len(n.pendingTXs) != 0 {
		t.Fatal("no pending TXs should be left to mine")
	}
}

func getTestDataDirPath() string {
	return filepath.Join(os.TempDir(), ".tbb_test")
}