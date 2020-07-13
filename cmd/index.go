package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use: "index",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Index......")
		log.Println(esClients)
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
