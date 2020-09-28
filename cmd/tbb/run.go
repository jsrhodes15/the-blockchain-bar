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
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ :=  cmd.Flags().GetUint64(flagPort)

			fmt.Printf("Launching TBB node and its HTTP API...\n\t- Configuration and data in %s directory.\n", getDataDirFromCmd(cmd))

			bootstrap := node.NewPeerNode(
				"127.0.0.1",
					8080,
					true,
					true,
				)

			n := node.New(getDataDirFromCmd(cmd), ip, port, bootstrap)
			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultFlags(runCmd)
	runCmd.Flags().String(flagIP, node.DefaultIP, "exposed IP address for communication with peers")
	runCmd.Flags().Uint64P(flagPort, "p", node.DefaultHttpPort, "exposed HTTP port for communication with peers")

	return runCmd
}
