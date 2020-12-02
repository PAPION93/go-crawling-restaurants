package domain

import (
	"time"
)

type Tabler interface {
	TableName() string
}

type Restaurant struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Point         string
	Address       string
	AddressDetail string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Restaurant) TableName() string {
	return "dining_restaurants"
}

type RestaurantUsecase interface {
	Create(*Restaurant) error
	Update(*Restaurant) error
}

type RestaurantRepository interface {
	GetRestaurant(*Restaurant) (Restaurant, error)
	Create(*Restaurant) error
	Update(*Restaurant) error
}
