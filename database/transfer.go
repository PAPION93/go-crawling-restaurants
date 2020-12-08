package database

import (
	"log"

	"tmwuw.com/domain"
)

// Transfer is transfer to real Table
type Transfer interface {
	TransferData()
}

type transfer struct {
	ru  domain.RestaurantUsecase
	rru domain.RestaurantRealUsecase
}

// NewTransfer ...
func NewTransfer(ru domain.RestaurantUsecase, rru domain.RestaurantRealUsecase) Transfer {
	return &transfer{
		ru,
		rru,
	}
}

func (t *transfer) TransferData() {
	page := 1
	for {
		restaurants, err := t.ru.GetLimit(page, 10)
		checkErr(err)
		if len(restaurants) == 0 {
			break
		}

		for _, restaurant := range restaurants {
			t.rru.CreateOrUpdate(&restaurant)
		}
		page++
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
