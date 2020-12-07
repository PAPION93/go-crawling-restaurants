package naver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"tmwuw.com/common"
	"tmwuw.com/domain"
)

// Naver interface
type Naver interface {
	CrawlLocation()
	CrawlGeocoding()
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

// Naver 지역검색 API를 통해 음식점의 주소와 카테고리를 가져온다.
func (n *naver) CrawlLocation() {
	page := 1
	for {
		restaurants, err := n.ru.GetLimit(page, 1)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}

		for _, restaurant := range restaurants {
			resp := requestLocation(restaurant.Address, restaurant.Name)
			if false == isWantedData(resp) {
				// time.Sleep(time.Second * 1)
				resp = requestLocation("대구", restaurant.Name)
				if false == isWantedData(resp) {
					// time.Sleep(time.Second * 1)
					resp = requestLocation("대구 "+restaurant.Address, restaurant.Name)
					if false == isWantedData(resp) {
						// time.Sleep(time.Second * 1)
						resp = requestLocation("", restaurant.Name)
						if false == isWantedData(resp) {
							time.Sleep(time.Millisecond * 500)
							continue
						}
					}
				}
			}

			err := n.ru.Update(&domain.Restaurant{ID: restaurant.ID, AddressDetail: resp.Items[0].Address, Category: resp.Items[0].Category})
			checkErr(err)
			time.Sleep(time.Millisecond * 500)
		}
		page++
		break
	}
}

// Naver Geocoding Api를 통해 위경도 정보를 가져온다.
func (n *naver) CrawlGeocoding() {

	page := 1
	for {
		restaurants, err := n.ru.GetLimit(page, 10)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}

		ch := make(chan common.ResponseLocation)
		for _, restaurant := range restaurants {
			go requestGeocoding(restaurant.Address, restaurant.AddressDetail, ch)
		}

		// 	err := n.ru.Update(&domain.Restaurant{ID: restaurant.ID, AddressDetail: resp.Items[0].Address, Category: resp.Items[0].Category})
		// 	checkErr(err)
		// 	time.Sleep(time.Millisecond * 500)
		page++
	}

}

func requestLocation(address string, name string) common.ResponseLocation {

	requestURI := createURI(address, name)
	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		checkErr(err)
	}

	req.Header.Add("X-Naver-Client-Id", viper.GetString(`secret.X-Naver-Client-Id`))
	req.Header.Add("X-Naver-Client-Secret", viper.GetString(`secret.X-Naver-Client-Secret`))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		checkErr(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	// str := string(bytes) //바이트를 문자열로
	// fmt.Println(str)

	var res common.ResponseLocation
	err = json.Unmarshal(bytes, &res)
	checkErr(err)

	return res
}

func requestGeocoding(address string, addressDetail string, ch chan<- common.ResponseLocation) {

	var requestURI string
	if addressDetail != "" {
		requestURI = "https://naveropenapi.apigw.ntruss.com/map-geocode/v2/geocode?query=" + addressDetail
	} else {
		requestURI = "https://naveropenapi.apigw.ntruss.com/map-geocode/v2/geocode?query=" + address
	}

	req, err := http.NewRequest("GET", url.QueryEscape(requestURI), nil)
	if err != nil {
		checkErr(err)
	}

	req.Header.Add("X-NCP-APIGW-API-KEY-ID", viper.GetString(`secret.X-NCP-APIGW-API-KEY-ID`))
	req.Header.Add("X-NCP-APIGW-API-KEY", viper.GetString(`secret.X-NCP-APIGW-API-KEY`))

}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func createURI(address string, name string) string {
	query := strings.Trim(address+" "+name, " ")
	escapedQuery := url.QueryEscape(query)
	log.Println("https://openapi.naver.com/v1/search/local?query=" + query)
	return "https://openapi.naver.com/v1/search/local?query=" + escapedQuery
}

func isWantedData(resp common.ResponseLocation) bool {
	if resp.Total == 0 {
		return false
	}

	if !strings.Contains(resp.Items[0].Address, "대구") {
		return false
	}

	return true
}
