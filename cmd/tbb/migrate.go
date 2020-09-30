package main

import (
	"context"
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"github.com/jsrhodes15/the-blockchain-bar/node"
	"github.com/spf13/cobra"
	"os"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk(getDataDirFromCmd(cmd))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			pendingBlock := node.NewPendingBlock(
				database.Hash{},
				state.NextBlockNumber(),
				[]database.Tx{
					database.NewTx("jrhodes", "jrhodes", 3, ""),
					database.NewTx("jrhodes", "jrhodes", 700, "reward"),
					database.NewTx("jrhodes", "meads", 2000, ""),
					database.NewTx("jrhodes", "jrhodes", 100, "reward"),
					database.NewTx("meads", "jrhodes", 1, ""),
					database.NewTx("meads", "lhendricks", 1000, ""),
					database.NewTx("meads", "jrhodes", 50, ""),
					database.NewTx("jrhodes", "jrhodes", 600, "reward"),
					database.NewTx("jrhodes", "jrhodes", 24700, "reward"),
				},
			)

			_, err = node.Mine(context.Background(), pendingBlock)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}

	addDefaultFlags(migrateCmd)

	return migrateCmd
}