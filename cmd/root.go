package cmd

import (
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Clusters map[string]Cluster
}

type Cluster struct {
	Addresses []string
}

var C Config

var esClients []*elastic.Client

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
	cobra.OnInitialize(initConfig, initClients)
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

func initClients() {
	for _, cluster := range C.Clusters {
		client, err := elastic.NewSimpleClient(elastic.SetURL(cluster.Addresses...))
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("connect to es fail")
			break
		}
		esClients = append(esClients, client)
	}
}
