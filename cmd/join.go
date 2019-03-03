package cmd

import (
	"fmt"

	"github.com/mikloslorinczi/snake-hub/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join to a Snake-hub server",
	Long: `
	Join to a Snake-hub server
	`,
	Run: func(cmd *cobra.Command, args []string) {
		client.Run()
	},
}

func init() {
	RootCmd.AddCommand(joinCmd)
	joinCmd.Flags().StringP("url", "u", "localhost:4545", "Snake-hub URL")
	if err := viper.BindPFlag("SNAKE_URL", joinCmd.Flags().Lookup("url")); err != nil {
		fmt.Printf("Cannot bind flag url %v\n", err)
	}

	joinCmd.Flags().StringP("name", "n", "", "Your name (max 8 character)")
	if err := viper.BindPFlag("SNAKE_USERNAME", joinCmd.Flags().Lookup("name")); err != nil {
		fmt.Printf("Cannot bind flag name %v\n", err)
	}

}
