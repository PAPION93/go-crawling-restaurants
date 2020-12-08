package naver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"tmwuw.com/common"
	"tmwuw.com/domain"
)

// Naver interface
type Naver interface {
	CrawlLocation()
	CrawlGeoLocation()
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
	page := 100
	for {
		restaurants, err := n.ru.GetLimit(page, 100)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}

		for _, restaurant := range restaurants {
			resp := requestLocation(restaurant.Address, restaurant.Name)
			if false == hasDaegu(resp) {
				resp = requestLocation("대구 "+restaurant.Address, restaurant.Name)
				if false == hasDaegu(resp) {
					resp = requestLocation("대구", restaurant.Name)
					if false == hasDaegu(resp) {
						time.Sleep(time.Millisecond * 200)
						continue
					}
				}
			}

			err := n.ru.Update(&domain.Restaurant{ID: restaurant.ID, AddressDetail: resp.Items[0].Address, Category: resp.Items[0].Category})
			checkErr(err)
			time.Sleep(time.Millisecond * 200)
		}
		page++
	}
}

// Naver GeoLocation Api를 통해 위경도 정보를 가져온다.
func (n *naver) CrawlGeoLocation() {

	page := 1
	for {
		restaurants, err := n.ru.GetLimit(page, 10)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}

		ch := make(chan common.ChannelResponseGeoLocation)
		for _, restaurant := range restaurants {
			go requestGeoLocation(restaurant.ID, restaurant.Address, restaurant.AddressDetail, ch)
		}

		for i := 0; i < len(restaurants); i++ {
			if resp, success := <-ch; success {
				if !hasLandNumber(resp) {
					continue
				}
				lat, _ := strconv.ParseFloat(resp.Addresses[0].X, 64)
				lng, _ := strconv.ParseFloat(resp.Addresses[0].Y, 64)
				err := n.ru.Update(&domain.Restaurant{ID: resp.RestaurantID, AddressDetail: resp.Addresses[0].JibunAddress, Lat: lat, Lng: lng})
				checkErr(err)
			} else {
				close(ch)
				break
			}
		}
		page++
		time.Sleep(time.Second * 1)
	}
}

func requestLocation(address string, name string) common.ResponseLocation {

	requestURI := setAddressAPIURI(address, name)
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

	var res common.ResponseLocation
	err = json.Unmarshal(bytes, &res)
	checkErr(err)

	return res
}

func requestGeoLocation(restaurantID uint, address string, addressDetail string, ch chan<- common.ChannelResponseGeoLocation) {

	var requestURI string

	if addressDetail != "" {
		requestURI = setGeoLocationAPIURI(addressDetail)
	} else {
		requestURI = setGeoLocationAPIURI(address)
	}

	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		checkErr(err)
	}

	req.Header.Add("X-NCP-APIGW-API-KEY-ID", viper.GetString(`secret.X-NCP-APIGW-API-KEY-ID`))
	req.Header.Add("X-NCP-APIGW-API-KEY", viper.GetString(`secret.X-NCP-APIGW-API-KEY`))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		checkErr(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	var res common.ResponseGeoLocation
	err = json.Unmarshal(bytes, &res)
	checkErr(err)

	ch <- common.ChannelResponseGeoLocation{RestaurantID: restaurantID, ResponseGeoLocation: res}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func setAddressAPIURI(address string, name string) string {
	escapedQuery := setQueryString(address, name)
	log.Println("setAddressAPIURI: "+address+" "+name, " ")
	return "https://openapi.naver.com/v1/search/local?query=" + escapedQuery
}

func setGeoLocationAPIURI(address string) string {
	escapedQuery := setQueryString(address, "")
	log.Println("setGeoLocationAPIURI: " + address)
	return "https://naveropenapi.apigw.ntruss.com/map-geocode/v2/geocode?query=" + escapedQuery
}

func setQueryString(address string, name string) (query string) {
	query = strings.Trim(address+" "+name, " ")
	query = url.QueryEscape(query)
	return
}

func hasDaegu(resp common.ResponseLocation) bool {
	if resp.Total == 0 {
		return false
	}

	if !strings.Contains(resp.Items[0].Address, "대구") {
		return false
	}

	return true
}

func hasLandNumber(resp common.ChannelResponseGeoLocation) bool {
	if resp.Status != "OK" {
		log.Println(resp)
		return false
	}

	if resp.Meta.TotalCount == 0 {
		return false
	}

	if resp.Addresses[0].AddressElements[7].LongName == "" {
		return false
	}
	return true
}
