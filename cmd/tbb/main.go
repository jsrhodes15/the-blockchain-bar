package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jsrhodes15/the-blockchain-bar/fs"
	"os"
	"path/filepath"
)

const flagDataDir = "datadir"
const flagMiner = "miner"
const flagIP = "ip"
const flagPort = "port"

const defaultDataDirname = ".tbb"

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(migrateCmd())
	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(runCmd())
	tbbCmd.AddCommand(balancesCmd())
	// TODO refactor txCmd so we can add it back to list of commands
	//tbbCmd.AddCommand(txCmd())

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		flagDataDir,
		"d",
		getDefaultDataDir(),
		"Absolute path to node data dir where DB is/will be stored",
	)
}

func getDataDirFromCmd(cmd *cobra.Command) string {
	dataDir, _ := cmd.Flags().GetString(flagDataDir)

	return fs.ExpandPath(dataDir)
}

func getDefaultDataDir() string {
	homeDir := fs.GetHomeDir()
	defaultDataDir := filepath.Join(homeDir, defaultDataDirname)

	return defaultDataDir
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
