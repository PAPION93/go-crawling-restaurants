package repository

import (
	"github.com/jinzhu/gorm"
	"tmwuw.com/domain"
)

// RestaurantRealRepository 구조체
type RestaurantRealRepository struct {
	DB *gorm.DB
}

// NewRestaurantRealRepository 생성자
func NewRestaurantRealRepository(DB *gorm.DB) domain.RestaurantRealRepository {
	return &RestaurantRealRepository{DB}
}

// GetRestaurant get By name and address
func (r *RestaurantRealRepository) GetRestaurant(name string, address string) (restaurant domain.RestaurantReal, err error) {
	restaurant = domain.RestaurantReal{}
	result := r.DB.Model(&domain.RestaurantReal{}).Where("name = ? AND address = ?", name, address).First(&restaurant)

	if result.Error != nil { // not found record
		return restaurant, result.Error
	}
	return restaurant, nil
}

// Create create
func (r *RestaurantRealRepository) Create(restaurant *domain.RestaurantReal) error {
	result := r.DB.Create(restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update update
func (r *RestaurantRealRepository) Update(restaurant *domain.RestaurantReal) error {
	result := r.DB.Model(&domain.RestaurantReal{}).Where("ID = ?", restaurant.ID).Updates(restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
