package node

import (
	"context"
	"encoding/hex"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"testing"
	"time"
)

func TestValidBlockHash(t *testing.T) {
	// Create a random hex string starting with 6 zeroes
	hexHash := "000000fa04f8160395c387277f8b2f14837603383d33809a4db586086168edfa"
	var hash = database.Hash{}
	// Convert it to raw bytes
	hex.Decode(hash[:], []byte(hexHash))
	// Validate the hash
	isValid := database.IsBlockHashValid(hash)
	if !isValid {
		t.Fatalf("hash '%s' with 6 zeroes should be valid", hexHash)
	}
}

func TestInvalidBlockHash(t *testing.T) {
	hexHash := "000001fa04f8160395c387277f8b2f14837603383d33809a4db586086168edfa"
	var hash = database.Hash{}

	hex.Decode(hash[:], []byte(hexHash))

	isValid := database.IsBlockHashValid(hash)
	if isValid {
		t.Fatal("hash is not suppose to be valid")
	}
}

func TestMine(t *testing.T) {
	pendingBlock := createRandomPendingBlock()

	ctx := context.Background()

	minedBlock, err := Mine(ctx, pendingBlock)
	if err != nil {
		t.Fatal(err)
	}

	minedBlockHash, err := minedBlock.Hash()
	if err != nil {
		t.Fatal(err)
	}

	if !database.IsBlockHashValid(minedBlockHash) {
		t.Fatal()
	}
}

func TestMineWithTimeout(t *testing.T) {
	pendingBlock := createRandomPendingBlock()

	ctx, _ := context.WithTimeout(context.Background(), time.Microsecond*100)

	_, err := Mine(ctx, pendingBlock)
	if err == nil {
		t.Fatal(err)
	}
}

func createRandomPendingBlock() PendingBlock {
	return NewPendingBlock(
		database.Hash{},
		0,
		[]database.Tx{
			database.NewTx("jrhodes", "jrhodes", 3, ""),
			database.NewTx("jrhodes", "jrhodes", 700, "reward"),
		},
	)
}