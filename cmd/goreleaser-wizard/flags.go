package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getBoolFlag reads a bool flag from the command.
// Flags are registered during init(), so a lookup error is a programming bug
// (misspelled literal name). The Must idiom makes it loud instead of silently
// reading the zero value, which would generate wrong configurations.
func getBoolFlag(cmd *cobra.Command, name string) bool {
	value, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(fmt.Sprintf("getBoolFlag: flag %q not defined (register it in init()): %v", name, err))
	}

	return value
}

// getStringFlag reads a string flag from the command.
// Flags are registered during init(), so a lookup error is a programming bug
// (misspelled literal name). The Must idiom makes it loud instead of silently
// reading the zero value, which would generate wrong configurations.
func getStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(fmt.Sprintf("getStringFlag: flag %q not defined (register it in init()): %v", name, err))
	}

	return value
}
