package diningcode

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"tmwuw.com/domain"
)

// Diningcode interface
type Diningcode interface {
	Crawl()
}

type diningcode struct {
	ru domain.RestaurantUsecase
}

// NewDiningcode ...
func NewDiningcode(ru domain.RestaurantUsecase) Diningcode {
	return &diningcode{
		ru,
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
func (d *diningcode) Crawl() {

	location := viper.GetStringMapString(`location`)
	for _, val := range location {
		for page := 1; page <= 10; page++ {

			// Request the HTML page.
			log.Printf("query: %s page: %d", val, page)
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

			if doc.Find("#div_list li").Length() == 0 {
				log.Fatal(doc.Find("body").Html())
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

					err := d.ru.Create(&domain.Restaurant{Name: name, Point: point, Address: address, AddressDetail: addressDetail})
					if err != nil {
						log.Fatal(err)
					}
				}
			})

			time.Sleep(time.Second * 5)
		}
	}
}
