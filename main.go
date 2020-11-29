package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"tmwuw.com/common"
	"tmwuw.com/diningcode"
)

func init() {
	viper.SetConfigFile(`./config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	dbConn := common.Init()
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	diningcodeRepo := diningcode.NewDiningcodeRepository(dbConn)
	diningcodeRepo.Crawl()
}
