package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Scrape()
}

// Scrape diningcode
func Scrape() {

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
		menu := s.Find(".stxt").Text()
		address := s.Find(".ctxt").Text()

		fmt.Printf("Name: %s, Menu: %s, Address: %s\n", name, menu, address)
	})
}
