package diningcode

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

// DiningcodeRepository interface
type DiningcodeRepository interface {
	Crawl() []restaurantInfo
}

type originalDiningcodeRepository struct {
	DB *sql.DB
}

// NewDiningcodeRepository ...
func NewDiningcodeRepository(db *sql.DB) DiningcodeRepository {
	return &originalDiningcodeRepository{
		DB: db,
	}
}

// Restaurant DataSet
type restaurantInfo struct {
	name          string
	point         string
	address       string
	addressDetail string
}

// Scrape Diningcode
func (r *originalDiningcodeRepository) Crawl() (restaurant []restaurantInfo) {

	subway := viper.GetStringMapString(`subway`)
	for _, val := range subway {
		for page := 1; page <= 10; page++ {

			// Request the HTML page.
			res, err := http.Get("https://www.diningcode.com/list.php?query=" + val + "&page=" + strconv.Itoa(page))
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

			// Find the restaurant items
			doc.Find("#div_list li").Each(func(i int, s *goquery.Selection) {
				name := s.Find(".btxt").Text()
				point := s.Find(".point").Text()

				s.Find(".ctxt").First().Remove()
				address := s.Find(".ctxt").Children().Remove().Text()
				addressDetail := s.Find(".ctxt").Text()

				if name != "" {
					// clear data
					index := strings.Index(name, ".")
					name = name[index+2:]

					point = strings.Split(point, "점")[0]

					restaurant = append(restaurant, restaurantInfo{name: name, point: point, address: address, addressDetail: addressDetail})
				}
			})
		}
	}

	return restaurant
}
