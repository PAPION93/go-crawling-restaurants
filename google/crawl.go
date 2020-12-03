package google

import (
	"log"
	"strings"
	"time"

	"github.com/fedesog/webdriver"
	"github.com/spf13/viper"
	"github.com/tebeka/selenium"
	"tmwuw.com/domain"
)

// Google interface
type Google interface {
	Crawl()
}

type google struct {
	ru domain.RestaurantUsecase
}

// NewGoogle ...
func NewGoogle(ru domain.RestaurantUsecase) Google {
	return &google{
		ru,
	}
}

// Crawl Google api
func (g *google) Crawl() {

	chromeDriver := webdriver.NewChromeDriver("/usr/local/bin/chromedriver")
	err := chromeDriver.Start()
	if err != nil {
		log.Println(err)
	}
	desired := webdriver.Capabilities{"Platform": "Mac"}
	required := webdriver.Capabilities{}

	locations := viper.GetStringMapString(`location`)
	for _, location := range locations {

		session, err := chromeDriver.NewSession(desired, required)
		if err != nil {
			log.Println(err)
		}

		err = session.Url("https://google.com")
		checkErr(err)
		// Enter code in textarea
		textInput, _ := session.FindElement(selenium.ByCSSSelector, ".gLFyf")
		textInput.Clear()
		textInput.SendKeys(location)
		textInput.Submit()

		// Search
		// btn, err := session.FindElement(selenium.ByCSSSelector, ".gNO89b")
		// checkErr(err)
		// btn.Click()

		// 더보기
		moreBtn, err := session.FindElement(selenium.ByCSSSelector, ".wUrVib")
		checkErr(err)
		moreBtn.Click()

		// err = session.FocusOnFrame("list")
		// if err != nil {
		// 	log.Println(err)
		// }

		for {

			restaurants, _ := session.FindElements(selenium.ByCSSSelector, ".rlfl__tls > div")
			for i, element := range restaurants {

				// Name 없을 경우 skip
				isPresent, err := element.FindElements(selenium.ByCSSSelector, ".dbg0pd > div")
				if len(isPresent) == 0 {
					continue
				}

				// Name
				nameElement, err := element.FindElement(selenium.ByCSSSelector, ".dbg0pd > div")
				checkErr(err)
				name, _ := nameElement.Text()

				// Point, Address
				details, err := element.FindElements(selenium.ByCSSSelector, ".rllt__details > div")
				var point, address string
				for j, detail := range details {
					if j == 0 {
						point, err = detail.Text()
						if point != "" {
							point = point[0:3]
						}
						checkErr(err)
					}

					if j == 1 {
						address, err = detail.Text()
						address = strings.Trim(address, " ")
						checkErr(err)
					}
				}
				log.Printf("%d:: name: %s, point: %s, address: %s\n", i, name, point, address)
			}

			isPresent, err := session.FindElements(selenium.ByCSSSelector, ".d6cvqb > a#pnnext")
			if len(isPresent) == 0 {
				log.Println(location + " 완료!!!")
				break
			}

			nextButton, err := session.FindElement(selenium.ByCSSSelector, ".d6cvqb > a#pnnext")
			checkErr(err)
			nextButton.Click()
			time.Sleep(10 * time.Second)
		}

		session.Delete()
		chromeDriver.Stop()
	}

}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
