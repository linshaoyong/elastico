package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/olivere/elastic"
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic("Fatal error config file: %s \n", err)
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
		log.Panic("Fatal error config file: %s \n", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Panic("Fatal error config file: %s \n", err)
	}
}

func initClients() {
	for _, cluster := range C.Clusters {
		client, err := elastic.NewSimpleClient(elastic.SetURL(cluster.Addresses...))
		if err != nil {
			log.Warn("connect to es fail: %s \n", err)
			break
		}
		esClients = append(esClients, client)
	}
}
