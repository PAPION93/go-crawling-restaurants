package repository

import (
	"github.com/jinzhu/gorm"
	"tmwuw.com/domain"
)

type RestaurantRepository struct {
	DB *gorm.DB
}

func NewRestaurantRepository(DB *gorm.DB) domain.RestaurantRepository {
	return &RestaurantRepository{DB}
}

func (r *RestaurantRepository) Create(restaurant *domain.Restaurant) error {
	result := r.DB.Create(restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RestaurantRepository) Update(restaurant *domain.Restaurant) error {
	result := r.DB.Model(&domain.Restaurant{}).Where("ID = ?", restaurant.ID).Save(&restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
