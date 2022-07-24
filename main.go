package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"

	"trader/services"
	_ "trader/services/stock"

	"trader/presentation/rest"
)

func init() {
	configPath := flag.String("config", "", "Service config file")

	if len(*configPath) == 0 {
		*configPath = "./"
	}

	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config from %s: %s", *configPath, err)
		os.Exit(1)
	}
}

func main() {
	var wg sync.WaitGroup

	rest.Run(&wg, services.ListAll())

	wg.Wait()
}
