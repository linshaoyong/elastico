package cmd

import (
	"context"
	"strings"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Clusters map[string]Cluster
}

// Cluster ...
type Cluster struct {
	Name          string
	ESClient      *elastic.Client
	Addresses     []string
	IndexOpenDays map[string]int64 `mapstructure:"index_open_days"`
}

// GetOpenedIndexNames ...
func (cluster *Cluster) GetOpenedIndexNames() []string {
	var opens []string
	names, err := cluster.ESClient.IndexNames()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("get opening indexes fail")
	}
	for _, name := range names {
		if !strings.HasPrefix(name, ".") {
			opens = append(opens, name)
		}
	}
	return opens
}

// GetClosedIndexNames ...
func (cluster *Cluster) GetClosedIndexNames() []string {
	var closeds []string
	ctx := context.Background()
	res, err := cluster.ESClient.CatIndices().Columns("index", "status").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range res {
		log.WithFields(log.Fields{"index": row.Index, "status": row.Status}).Warn("index")
		if row.Status == "close" && !strings.HasPrefix(row.Index, ".") {
			closeds = append(closeds, row.Index)
		}
	}
	return closeds
}

// CloseIndex ...
func (cluster *Cluster) CloseIndex(name string) {
	cresp, err := cluster.ESClient.CloseIndex(name).Do(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name, "error": err}).Warn("close index fail")
	} else if !cresp.Acknowledged {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name, "error": err}).Warn("expected close index to be acknowledged")
	} else {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name}).Info("old index closed")
	}
}

// C ...
var C Config

var clusters []Cluster

var rootCmd = &cobra.Command{
	Use: "elastico",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Root......")
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"error": err}).Panic("Fatal error config file")
	}
}

func init() {
	cobra.OnInitialize(initConfig, initClusters)
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("elastico")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Panic("Read config file fail")
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Panic("Unmarshal config file fail")
	}
}

func initClusters() {
	for name, cluster := range C.Clusters {
		client, err := elastic.NewSimpleClient(elastic.SetURL(cluster.Addresses...))
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("connect to es fail")
			break
		}
		cluster.Name = name
		cluster.ESClient = client
		clusters = append(clusters, cluster)
	}
}
