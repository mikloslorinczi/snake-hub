package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "snake-hub",
	Short: "snake-hub is a multiplayer snake game",
	Long: `
Start a server with snake-hub host
users can connect with snake-hub join
`,
}
