package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikloslorinczi/snake-hub/server"
)

// hostCmd represents the run command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Starts the Snake-hub server",
	Long: `
	Starts the Snake-hub server
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SNAKE_PORT set to %v\n", viper.Get("SNAKE_PORT"))
		fmt.Printf("SNAKE_SECRET set to %v\n", viper.Get("SNAKE_SECRET"))
		// server.Setup()
		server.Run()
	},
}

func init() {

	RootCmd.AddCommand(hostCmd)

	hostCmd.Flags().IntP("port", "p", 4545, "Snake-hub port")
	if err := viper.BindPFlag("SNAKE_PORT", hostCmd.Flags().Lookup("port")); err != nil {
		fmt.Printf("Cannot bind flag SNAKE_PORT %v\n", err)
	}

}
