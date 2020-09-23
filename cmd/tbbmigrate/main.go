package main

import (
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"os"
	"time"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("jrhodes", "jrhodes", 3, ""),
			database.NewTx("jrhodes", "jrhodes", 700, "reward"),
		},
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("jrhodes", "meads", 2000, ""),
			database.NewTx("jrhodes", "jrhodes", 100, "reward"),
			database.NewTx("meads", "jrhodes", 1, ""),
			database.NewTx("meads", "lhendricks", 1000, ""),
			database.NewTx("meads", "jrhodes", 50, ""),
			database.NewTx("jrhodes", "jrhodes", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}