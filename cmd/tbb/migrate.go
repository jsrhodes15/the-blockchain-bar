package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"os"
	"time"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			block0 := database.NewBlock(
				database.Hash{},
				0,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("jrhodes", "jrhodes", 3, ""),
					database.NewTx("jrhodes", "jrhodes", 700, "reward"),
				},
			)

			state.AddBlock(block0)
			block0hash, err := state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block1 := database.NewBlock(
				block0hash,
				1,
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
			block1hash, err := state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			block2 := database.NewBlock(
				block1hash,
				2,
				uint64(time.Now().Unix()),
				[]database.Tx{
					database.NewTx("jrhodes", "jrhodes", 24700, "reward"),
				},
			)

			state.AddBlock(block2)
			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	addDefaultFlags(migrateCmd)

	return migrateCmd
}
