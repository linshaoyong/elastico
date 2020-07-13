package main

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v6"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Clusters map[string]Cluster
}

type Cluster struct {
	Addresses []string
}

var C Config

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("conf")
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

func main() {
	for _, cluster := range C.Clusters {
		cfg := elasticsearch.Config{
			Addresses: cluster.Addresses,
		}
		es, err := elasticsearch.NewClient(cfg)
		if err != nil {
			log.Warn("connect to es fail: %s \n", err)
			break
		}
		log.Println(es.Info())
	}
}
