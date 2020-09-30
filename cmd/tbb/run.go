package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jsrhodes15/the-blockchain-bar/database"
	"github.com/jsrhodes15/the-blockchain-bar/node"
	"os"
)

func runCmd() *cobra.Command{
	var runCmd = &cobra.Command{
		Use: "run",
		Short: "Launches the TBB node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			miner, _ := cmd.Flags().GetString(flagMiner)
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ := cmd.Flags().GetUint64(flagPort)

			fmt.Printf("Launching TBB node and its HTTP API...\n\t- Configuration and data in %s directory.\n", getDataDirFromCmd(cmd))

			bootstrap := node.NewPeerNode(
				"127.0.0.1",
					8080,
					true,
					database.NewAccount("jrhodes"),
					false,
				)

			n := node.New(getDataDirFromCmd(cmd), ip, port, database.NewAccount(miner), bootstrap)
			err := n.Run(context.Background())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultFlags(runCmd)
	runCmd.Flags().StringP(flagMiner, "m", node.DefaultMiner, "miner account of this node to receive block rewards")
	runCmd.Flags().String(flagIP, node.DefaultIP, "exposed IP address for communication with peers")
	runCmd.Flags().Uint64P(flagPort, "p", node.DefaultHttpPort, "exposed HTTP port for communication with peers")

	return runCmd
}
