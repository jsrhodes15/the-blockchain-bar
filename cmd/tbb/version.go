package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "9"
const Patch = "0"
const Description = "Proof of Work Reward"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s-beta %s", Major, Minor, Patch, Description))
	},
}
