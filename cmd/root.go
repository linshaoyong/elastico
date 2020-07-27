package cmd

import (
	"context"
	"strings"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Clusters map[string]Cluster
}

type Cluster struct {
	Name                   string
	ESClient               *elastic.Client
	Addresses              []string
	IndexOpenDays          map[string]int64 `mapstructure:"index_open_days"`
	IndexDefaultRemainDays int64            `mapstructure:"index_default_remain_days"`
	IndexDefaultOpenDays   int64            `mapstructure:"index_default_open_days"`
}

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

func (cluster *Cluster) GetClosedIndexNames() []string {
	var closeds []string
	ctx := context.Background()
	res, err := cluster.ESClient.CatIndices().Columns("index", "status").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range res {
		if row.Status == "close" && !strings.HasPrefix(row.Index, ".") {
			closeds = append(closeds, row.Index)
		}
	}
	return closeds
}

func (cluster *Cluster) CloseIndex(name string) {
	cresp, err := cluster.ESClient.CloseIndex(name).Do(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name, "error": err}).Warn("close index fail")
	} else if !cresp.Acknowledged {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name}).Warn("expected close index to be acknowledged")
	} else {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name}).Info("index closed")
	}
}

func (cluster *Cluster) DeleteIndex(name string) {
	ctx := context.Background()
	deleteIndex, err := cluster.ESClient.DeleteIndex(name).Do(ctx)
	if err != nil {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name, "error": err}).Warn("delete index fail")
	}
	if !deleteIndex.Acknowledged {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name}).Warn("expected delete index to be acknowledged")
	} else {
		log.WithFields(log.Fields{"cluster": cluster.Name, "index": name}).Info("index deleted")
	}
}

var C Config

var clusters []Cluster

var rootCmd = &cobra.Command{
	Use: "elastico",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Root......")
	},
}

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
			continue
		}
		cluster.Name = name
		cluster.ESClient = client
		if cluster.IndexDefaultRemainDays < 10 {
			cluster.IndexDefaultRemainDays = 125
		}
		clusters = append(clusters, cluster)
	}
}
