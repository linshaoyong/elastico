package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "index [list/close/delete]",
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := cmd.Flags().GetString("a")
		if a == "list" {
			list()
		} else if a == "close" {
			log.Println(a)
		} else if a == "delete" {
			log.Println(a)
		} else {
			log.Println("Do nothing......")
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.Flags().String("a", "list", "index action")
}

func list() {
	for _, client := range esClients {
		names, err := client.IndexNames()
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, name := range names {
			if !strings.HasPrefix(name, ".") {
				log.Println(name)
			}
		}
	}
}

func close() {

}
