package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"tmwuw.com/database"
	"tmwuw.com/domain"
	"tmwuw.com/restaurant/repository"
	"tmwuw.com/restaurant/usecase"
	r_repo "tmwuw.com/restaurant_real/repository"
	r_usecase "tmwuw.com/restaurant_real/usecase"
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
	a.DBInit()
	// a.TestDBInit()
	defer a.DBClose(a.DB)
	a.DB.AutoMigrate(&domain.Restaurant{})

	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	// 파일로그로 변경
	log.SetOutput(fpLog)
	log.Println("Start")

	rr := repository.NewRestaurantRepository(a.DB)
	ur := usecase.NewRestaurantUsecase(rr)

	// googleRepo := google.NewGoogle(ur)
	// googleRepo.Crawl()

	// diningcodeRepo := diningcode.NewDiningcode(ur)
	// diningcodeRepo.Crawl()

	// naverRepo := naver.NewNaver(ur)
	// naverRepo.CrawlLocation()
	// naverRepo.CrawlGeoLocation()

	rru := r_repo.NewRestaurantRealRepository(a.DB)
	rur := r_usecase.NewRestaurantRealUsecase(rru)
	tranfer := database.NewTransfer(ur, rur)
	tranfer.TransferData()

	log.Println("End")
}
