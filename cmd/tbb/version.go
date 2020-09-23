package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "4"
const Patch = "0"
const Description = "Blockchain"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s-beta %s", Major, Minor, Patch, Description)
	},
}
