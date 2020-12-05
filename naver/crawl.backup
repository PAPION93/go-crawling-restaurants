package naver

import (
	"fmt"
	"log"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/tebeka/selenium"
	"tmwuw.com/domain"
)

// Naver interface
type Naver interface {
	Crawl()
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

// Crawl Naver api
func (d *naver) Crawl() {

	chromeDriver := webdriver.NewChromeDriver("/usr/local/bin/chromedriver")
	err := chromeDriver.Start()
	if err != nil {
		log.Println(err)
	}
	desired := webdriver.Capabilities{"Platform": "Mac"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Println(err)
	}

	err = session.Url("https://naver.com")
	if err != nil {
		log.Println(err)
	}

	// Enter code in textarea
	elem, _ := session.FindElement(selenium.ByCSSSelector, ".green_window > #query")
	elem.Clear()
	elem.SendKeys("만촌동 맛집")

	// Click the run button
	btn, _ := session.FindElement(selenium.ByCSSSelector, "#search_btn")
	btn.Click()

	api_more, _ := session.FindElement(selenium.ByCSSSelector, ".api_more")
	api_more.Click()

	element, err := session.FindElement(selenium.ByCSSSelector, "#searchIframe")
	// Different ways to get element
	//element, err := wd.FindElement(selenium.ByCSSSelector, `iframe[name="innerFrame"]`)
	if err != nil {
		panic(err)
	}

	//Switch to iframe
	err = session.FocusOnFrame(element)
	if err != nil {
		panic(err)
	}

	// err = session.FocusOnFrame("#searchIframe")
	// if err != nil {
	// 	log.Println(err)
	// }

	restaurants, err := session.FindElement(selenium.ByCSSSelector, "#_pcmap_list_scroll_container")
	fmt.Println(restaurants)
	// for val := range restaurants {
	// fmt.Println(val)
	// }

	time.Sleep(5 * time.Second)
	session.Delete()
	chromeDriver.Stop()

}
