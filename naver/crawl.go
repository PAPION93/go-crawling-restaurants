package naver

import (
	"fmt"

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
func (n *naver) Crawl() {
	restaurants, err := n.ru.GetLimit(2, 100)
	if err != nil {
		fmt.Println(err)
	}
	for key, val := range restaurants {
		fmt.Println(key, val)
	}

}
