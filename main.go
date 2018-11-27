package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/snake-hub/cmd"
	"github.com/mikloslorinczi/snake-hub/utils"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("Error during execution : %v\n", err)
	}
}

func init() {

	// Load config from snake-hub.yaml and ENV
	if err := utils.ReadConfig("./", "snake-hub", nil); err != nil {
		fmt.Printf("Cannot set configuration %v\n", err)
	}

	// Setup global flags
	cmd.RootCmd.PersistentFlags().StringP("secret", "s", "", "Snake-hub secret")

	if err := viper.BindPFlag("SNAKE_SECRET", cmd.RootCmd.PersistentFlags().Lookup("secret")); err != nil {
		fmt.Printf("Cannot bind flag SNAKE_SECRET %v\n", err)
	}

}
