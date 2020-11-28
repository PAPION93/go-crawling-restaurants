package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Restaurant struct {
	name           string
	point          string
	address        string
	address_detail string
}

func init() {
	viper.SetConfigFile(`./config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`postgres`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	restaurants := Scrape()
	for key, val := range restaurants {
		fmt.Println(key, val)
	}
}

// Scrape Diningcode
func Scrape() (restaurant []Restaurant) {

	query := "만촌역"
	page := 1

	// Request the HTML page.
	res, err := http.Get("https://www.diningcode.com/list.php?query=" + query + "&page=" + strconv.Itoa(page))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("#div_list li").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".btxt").Text()
		point := s.Find(".point").Text()

		s.Find(".ctxt").First().Remove()
		address := s.Find(".ctxt").Children().Remove().Text()
		address_detail := s.Find(".ctxt").Text()

		if name != "" {
			// clear data
			index := strings.Index(name, ".")
			name = name[index+2:]

			point = strings.Split(point, "점")[0]

			restaurant = append(restaurant, Restaurant{name: name, point: point, address: address, address_detail: address_detail})
		}
	})

	return restaurant
}
