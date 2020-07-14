package cmd

import (
	"strings"
	"time"

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

func getOpeningIndexNames() []string {
	var opens []string
	for _, client := range esClients {
		names, err := client.IndexNames()
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, name := range names {
			if !strings.HasPrefix(name, ".") {
				opens = append(opens, name)
			}
		}
	}
	return opens
}

func isOldIndex(name string, formats []string, days int64, now int64) bool {
	for _, format := range formats {
		if len(name) < len(format) {
			continue
		}
		ts, err := time.Parse(format, name[len(name)-len(format):len(name)])
		if err == nil {
			if now-ts.Unix() > 86400*days {
				return true
			}
		}
	}
	return false
}

func getOldIndexNames(indexNames []string, days int64) []string {
	var olds []string
	now := time.Now().Unix()
	for _, name := range indexNames {
		if isOldIndex(name, []string{"20060102", "2006.01.02"}, days, now) {
			olds = append(olds, name)
		}
	}
	return olds
}

func list() {
	for _, name := range getOpeningIndexNames() {
		log.Println(name)
	}
}

func close() {

}
