package main

import (
	"context"
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"github.com/jsrhodes15/the-blockchain-bar/node"
	"github.com/spf13/cobra"
	"time"
)

var migrateCmd = func() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrates the blockchain database according to new business rules.",
		Run: func(cmd *cobra.Command, args []string) {
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ := cmd.Flags().GetUint64(flagPort)

			peer := node.NewPeerNode(
				"127.0.0.1",
				8080,
				true,
				false,
				)

			n := node.New(getDataDirFromCmd(cmd), ip, port, peer)

			n.AddPendingTX(database.NewTx("jrhodes", "jrhodes", 3, ""), peer)
			n.AddPendingTX(database.NewTx("jrhodes", "meads", 2000, ""), peer)
			n.AddPendingTX(database.NewTx("meads", "jrhodes", 1, ""), peer)
			n.AddPendingTX(database.NewTx("meads", "lhendricks", 1000, ""), peer)
			n.AddPendingTX(database.NewTx("meads", "jrhodes", 50, ""), peer)

			ctx, closeNode := context.WithTimeout(context.Background(), time.Minute*15)

			go func() {
				ticker := time.NewTicker(time.Second * 10)

				for {
					select {
					case <-ticker.C:
						if !n.LatestBlockHash().IsEmpty() {
							closeNode()
							return
						}
					}
				}
			}()

			err := n.Run(ctx)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	addDefaultFlags(migrateCmd)
	migrateCmd.Flags().String(flagIP, node.DefaultIP, "exposed IP for communication with peers")
	migrateCmd.Flags().Uint64P(flagPort, "p", node.DefaultHttpPort, "exposed HTTP port for communication with peers")

	return migrateCmd
}