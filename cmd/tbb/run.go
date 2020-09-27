package main

import (
	"fmt"
	"github.com/jsrhodes15/the-blockchain-bar/node"
	"github.com/spf13/cobra"
	"os"
)

func runCmd() *cobra.Command{
	var runCmd = &cobra.Command{
		Use: "run",
		Short: "Launches the TBB node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir := getDataDirFromCmd(cmd)
			port, _ :=  cmd.Flags().GetUint64(flagPort)

			fmt.Println("Launching TBB node and its HTTP API...")
			fmt.Printf("Configuration and data in %s directory.\n\n", dataDir)

			bootstrap := node.NewPeerNode(
				"127.0.0.1",
					8080,
					true,
					true,
				)

			n := node.New(getDataDirFromCmd(cmd), port, bootstrap)
			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultFlags(runCmd)
	runCmd.Flags().Uint64P(flagPort, "p", node.DefaultHttpPort, "exposed HTTP port for communication with peers")
	return runCmd
}
