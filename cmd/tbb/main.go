package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const flagDataDir = "datadir"
const defaultDataDirname = ".tbb"
var defaultDataDir string

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defaultDataDir = filepath.Join(homeDir, defaultDataDirname)

	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(txCmd())
	tbbCmd.AddCommand(runCmd())

	err = tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		flagDataDir,
		"d",
		defaultDataDir,
		"Absolute path where tbb data is stored",
	)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
