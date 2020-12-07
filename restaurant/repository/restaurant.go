package repository

import (
	"github.com/jinzhu/gorm"
	"tmwuw.com/domain"
)

// RestaurantRepository 구조체
type RestaurantRepository struct {
	DB *gorm.DB
}

// NewRestaurantRepository 생성자
func NewRestaurantRepository(DB *gorm.DB) domain.RestaurantRepository {
	return &RestaurantRepository{DB}
}

// GetRestaurant get By name and address
func (r *RestaurantRepository) GetRestaurant(name string, address string) (domain.Restaurant, error) {
	result := r.DB.Model(&domain.Restaurant{}).Where("name = ? AND address = ?", name, address).First(domain.Restaurant{})
	if result.Error != nil {
		return domain.Restaurant{}, result.Error
	}
	return domain.Restaurant{}, nil
}

// GetLimit get restaurant list with size
func (r *RestaurantRepository) GetLimit(offset int, size int) ([]domain.Restaurant, error) {
	restaurants := []domain.Restaurant{}

	result := r.DB.Offset(offset).Limit(size).Order("id").Find(&restaurants)
	if result.Error != nil {
		return restaurants, result.Error
	}
	return restaurants, nil
}

// Create create
func (r *RestaurantRepository) Create(restaurant *domain.Restaurant) error {
	result := r.DB.Create(restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update update
func (r *RestaurantRepository) Update(restaurant *domain.Restaurant) error {
	result := r.DB.Model(&domain.Restaurant{}).Select("Category", "AddressDetail", "UpdatedAt").Where("ID = ?", restaurant.ID).Save(&restaurant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
