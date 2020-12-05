package naver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tmwuw.com/domain"
)

// Naver interface
type Naver interface {
	RequestAPI()
}

type naver struct {
	ru domain.RestaurantUsecase
}

// NewNaver ...
func NewNaver(ru domain.RestaurantUsecase) Naver {
	return &naver{
		ru,
	}
}

type responseData struct {
	Total int
	Items []struct {
		Title    string
		Link     string
		Category string
		Address  string
	}
}

// Request Naver api
func (n *naver) RequestAPI() {
	page := 1
	for {
		restaurants, err := n.ru.GetLimit(page, 100)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}
		for _, restaurant := range restaurants {
			res := requestAddress(restaurant.Address, restaurant.Name)
			time.Sleep(time.Second * 1)

			if res.Total == 0 {
				res = requestAddress("대구", restaurant.Name)
				time.Sleep(time.Second * 1)
				if res.Total == 0 {
					res = requestAddress("", restaurant.Name)
					time.Sleep(time.Second * 1)
					if res.Total == 0 {
						continue
					}
				}
			}

			err := n.ru.Update(&domain.Restaurant{ID: restaurant.ID, AddressDetail: res.Items[0].Address, Category: res.Items[0].Category})
			checkErr(err)

			time.Sleep(time.Second * 1)
		}
		page += 1
	}
}

func requestAddress(address string, name string) responseData {
	query := url.QueryEscape(strings.Trim(address+" "+name, " "))
	requestURI := "https://openapi.naver.com/v1/search/local?query=" + query
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		checkErr(err)
	}
	fmt.Println(strings.Trim(address+" "+name, " "))

	req.Header.Add("X-Naver-Client-Id", "sOCY0oNpt0S8nLDN38Wp")
	req.Header.Add("X-Naver-Client-Secret", "GtW_LGRrE3")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		checkErr(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	// str := string(bytes) //바이트를 문자열로
	// fmt.Println(str)

	var res responseData
	err = json.Unmarshal(bytes, &res)
	checkErr(err)

	return res
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
