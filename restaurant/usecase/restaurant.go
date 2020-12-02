package usecase

import (
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

func (r *restaurantUsecase) Create(restaurant *domain.Restaurant) error {
	_, err := r.restaurantRepo.GetRestaurant(restaurant)
	// if errors.Is(err.Error, gorm.ErrRecordNotFound) {
	if err != nil {
		return r.restaurantRepo.Create(restaurant)
	}
	return nil
}

func (r *restaurantUsecase) Update(restaurant *domain.Restaurant) error {
	return r.restaurantRepo.Update(restaurant)
}
