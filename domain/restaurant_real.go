package domain

import (
	"time"
)

type TablerR interface {
	TableName() string
}

type RestaurantReal struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Category    string
	Address     string
	GooglePoint string
	NaverPoint  string
	DiningPoint string
	Lat         float64
	Lng         float64
}

func (RestaurantReal) TableName() string {
	return "restaurants"
}

type RestaurantRealUsecase interface {
	CreateOrUpdate(*Restaurant) (RestaurantReal, error)
	Create(*RestaurantReal) error
	Update(*RestaurantReal) error
}

type RestaurantRealRepository interface {
	GetRestaurant(name string, address string) (RestaurantReal, error)
	Create(*RestaurantReal) error
	Update(*RestaurantReal) error
}
