package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "snake-hub",
	Short: "snake-hub is a Cobra app",
	Long: `
Long description of snake-hub
can span multiple
lines
`,
}
