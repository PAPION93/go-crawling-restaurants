package domain

import (
	"time"
)

type Tabler interface {
	TableName() string
}

type Restaurant struct {
	ID            uint `gorm:"primaryKey"`
	Category      string
	Name          string
	GooglePoint   string
	DiningPoint   string
	NaverPoint    string
	Address       string
	AddressDetail string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Restaurant) TableName() string {
	return "crwaling_restaurants"
}

type RestaurantUsecase interface {
	GetLimit(page int, size int) ([]Restaurant, error)
	Create(*Restaurant) error
	Update(*Restaurant) error
}

type RestaurantRepository interface {
	GetRestaurant(name string, address string) (Restaurant, error)
	GetLimit(offset int, size int) ([]Restaurant, error)
	Create(*Restaurant) error
	Update(*Restaurant) error
}
