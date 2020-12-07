package usecase

import (
	"github.com/jinzhu/gorm"
	"tmwuw.com/domain"
)

// RestaurantUsecase ...
type restaurantUsecase struct {
	restaurantRepo domain.RestaurantRepository
}

// NewRestaurantUsecase will create new an RestaurantUsecase object representation of domain.RestaurantUsecase interface
func NewRestaurantUsecase(r domain.RestaurantRepository) domain.RestaurantUsecase {
	return &restaurantUsecase{
		restaurantRepo: r,
	}
}

func (r *restaurantUsecase) GetLimit(page int, size int) ([]domain.Restaurant, error) {
	if page == 0 {
		page = 1
	}

	switch {
	case size > 100:
		size = 100
	case size <= 0:
		size = 10
	}

	offset := (page - 1) * size

	restaurants, err := r.restaurantRepo.GetLimit(offset, size)
	if err != nil {
		return restaurants, err
	}
	return restaurants, nil
}

func (r *restaurantUsecase) Create(restaurant *domain.Restaurant) error {
	_, err := r.restaurantRepo.GetRestaurant(restaurant.Name, restaurant.Address)
	// if errors.Is(err.Error, gorm.ErrRecordNotFound) {
	if err != nil {
		return r.restaurantRepo.Create(restaurant)
	}
	return nil
}

func (r *restaurantUsecase) Update(restaurant *domain.Restaurant) error {
	if restaurant.ID < 0 {
		return gorm.ErrRecordNotFound
	}

	return r.restaurantRepo.Update(restaurant)
}
