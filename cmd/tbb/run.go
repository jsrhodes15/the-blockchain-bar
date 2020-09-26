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
			dataDir, _ := cmd.Flags().GetString(flagDataDir)

			fmt.Println("Launching TBB node and its HTTP API...")
			fmt.Printf("Configuration and data in %s directory.\n\n", dataDir)

			err := node.Run(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(runCmd)

	return runCmd
}
