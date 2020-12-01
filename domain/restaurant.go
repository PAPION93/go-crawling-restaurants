package domain

import (
	"time"
)

type Restaurant struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Point         string
	Address       string
	AddressDetail string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type RestaurantUsecase interface {
	Create(*Restaurant) error
	Update(*Restaurant) error
}

type RestaurantRepository interface {
	Create(*Restaurant) error
	Update(*Restaurant) error
}
