package cmd

import (
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

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
			delete()
		} else {
			log.Println("Do nothing......")
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.Flags().String("a", "list", "index action")
}

func isIndexEarlyThan(name string, formats []string, startUnix int64) (bool, bool) {
	match := false
	for _, format := range formats {
		if len(name) < len(format) {
			continue
		}
		ts, err := time.Parse(format, name[len(name)-len(format):len(name)])
		if err == nil {
			match = true
			if startUnix > ts.Unix() {
				return true, match
			}
		}
	}
	return false, match
}

func filterIndexsEarlyThan(indexNames []string, days int64, customDays map[string]int64) []string {
	var olds []string
	now := time.Now().Unix()
	for _, name := range indexNames {
		var deltaDays = days
		for cdk, cdv := range customDays {
			if strings.HasPrefix(name, cdk) {
				deltaDays = cdv
				break
			}
		}

		startUnix := now - 86400*deltaDays
		// daily index
		isEarly, match := isIndexEarlyThan(name, []string{"20060102", "2006.01.02"}, startUnix)
		if isEarly {
			olds = append(olds, name)
		}

		if !match {
			// maybe monthly index
			deltaDays += 50
			startUnix = now - 86400*deltaDays
			if isEarly, _ = isIndexEarlyThan(name, []string{"2006.01", "200601"}, startUnix); isEarly {
				olds = append(olds, name)
			}
		}
	}
	return olds
}

func list() {
	for _, cluster := range clusters {
		for _, name := range cluster.GetOpenedIndexNames() {
			log.Info(name)
		}
	}
}

func close() {
	for _, cluster := range clusters {
		wg.Add(1)
		go func(c Cluster) {
			names := filterIndexsEarlyThan(c.GetOpenedIndexNames(), c.IndexDefaultOpenDays, c.IndexOpenDays)
			for _, name := range names {
				c.CloseIndex(name)
			}
			defer wg.Done()
		}(cluster)
	}
	wg.Wait()
}

func delete() {
	for _, cluster := range clusters {
		wg.Add(1)
		go func(c Cluster) {
			names := c.GetClosedIndexNames()
			names = filterIndexsEarlyThan(names, c.IndexDefaultRemainDays, c.IndexRemainDays)
			for _, name := range names {
				c.DeleteIndex(name)
			}
			defer wg.Done()
		}(cluster)
	}
	wg.Wait()
}
