package usecase

import (
	"github.com/jinzhu/gorm"
	"tmwuw.com/domain"
)

// RestaurantRealUsecase ...
type restaurantRealUsecase struct {
	restaurantRealRepo domain.RestaurantRealRepository
}

// NewRestaurantRealUsecase will create new an RestaurantRealUsecase object representation of domain.RestaurantRealUsecase interface
func NewRestaurantRealUsecase(r domain.RestaurantRealRepository) domain.RestaurantRealUsecase {
	return &restaurantRealUsecase{
		restaurantRealRepo: r,
	}
}

func (r *restaurantRealUsecase) CreateOrUpdate(restaurant *domain.Restaurant) (domain.RestaurantReal, error) {
	result, err := r.restaurantRealRepo.GetRestaurant(restaurant.Name, restaurant.AddressDetail)
	if err != nil {
		r.restaurantRealRepo.Create(&domain.RestaurantReal{
			Name:        restaurant.Name,
			Category:    restaurant.Category,
			Address:     restaurant.AddressDetail,
			GooglePoint: restaurant.GooglePoint,
			NaverPoint:  restaurant.NaverPoint,
			DiningPoint: restaurant.DiningPoint,
			Lat:         restaurant.Lat,
			Lng:         restaurant.Lng,
		})
	} else {
		// log.Println(result)
		r.restaurantRealRepo.Update(&domain.RestaurantReal{
			ID:          result.ID,
			Name:        restaurant.Name,
			Category:    restaurant.Category,
			Address:     restaurant.AddressDetail,
			GooglePoint: restaurant.GooglePoint,
			NaverPoint:  restaurant.NaverPoint,
			DiningPoint: restaurant.DiningPoint,
			Lat:         restaurant.Lat,
			Lng:         restaurant.Lng,
		})
	}
	return result, nil
}

func (r *restaurantRealUsecase) Create(restaurant *domain.RestaurantReal) error {
	_, err := r.restaurantRealRepo.GetRestaurant(restaurant.Name, restaurant.Address)
	if err != nil {
		return r.restaurantRealRepo.Create(restaurant)
	}
	return nil
}

func (r *restaurantRealUsecase) Update(restaurant *domain.RestaurantReal) error {
	if restaurant.ID < 1 {
		return gorm.ErrRecordNotFound
	}

	return r.restaurantRealRepo.Update(restaurant)
}
