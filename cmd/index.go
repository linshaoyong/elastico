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
			close()
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

func getOldIndexNames(indexNames []string, days int64, customDays map[string]int64) []string {
	var olds []string
	now := time.Now().Unix()
	for _, name := range indexNames {
		var realDays = days
		for cdk, cdv := range customDays {
			if strings.HasPrefix(name, cdk) {
				realDays = cdv
				break
			}
		}
		if isOldIndex(name, []string{"20060102", "2006.01.02"}, realDays, now) {
			olds = append(olds, name)
		}
	}
	return olds
}

func list() {
	for _, cluster := range clusters {
		for _, name := range cluster.GetOpeningIndexNames() {
			log.Info(name)
		}
	}
}

func close() {
	for _, cluster := range clusters {
		names := getOldIndexNames(cluster.GetOpeningIndexNames(), 7, cluster.IndexOpenDays)
		for _, name := range names {
			cluster.CloseIndex(name)
		}
	}
}
