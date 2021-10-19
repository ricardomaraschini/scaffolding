package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version holds the current binary version. Set at compile time.
var Version = "v0.0.0"

var version = &cobra.Command{
	Use:   "version",
	Short: "Prints the version",
	RunE: func(c *cobra.Command, args []string) error {
		fmt.Println("version", Version)
		return nil
	},
}
