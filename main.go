package main

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"tmwuw.com/database"
	"tmwuw.com/diningcode"
	"tmwuw.com/domain"
	"tmwuw.com/restaurant/repository"
	"tmwuw.com/restaurant/usecase"
)

func init() {
	viper.SetConfigFile(`./config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	a := database.Config{}
	// a.DBInit()
	a.TestDBInit()
	defer a.DBClose(a.DB)
	a.DB.AutoMigrate(&domain.Restaurant{})

	rr := repository.NewRestaurantRepository(a.DB)
	ur := usecase.NewRestaurantUsecase(rr)
	diningcodeRepo := diningcode.NewDiningcode(ur)
	diningcodeRepo.Crawl()
}
