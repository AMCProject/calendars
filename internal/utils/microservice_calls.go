package utils

import (
	"calendar/internal"
	"calendar/internal/config"
	"calendar/internal/models"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type Endpoints struct {
}
type EndpointsI interface {
	GetAllMeals(userId string) (meals []*models.MealToFront, err error)
	GetMeal(userId, mealId string) (meal models.MealToFront, err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetAllMeals(userId string) (meals []*models.MealToFront, err error) {
	url := config.Config.MealsURL + "user/" + userId + "/meal"
	season := getSeason()
	if season != "" {
		url += "?season[]=" + season
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return []*models.MealToFront{}, internal.ErrReturningAllMeals
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		newError := new(internal.ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return []*models.MealToFront{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&meals)
	if err != nil {
		log.Error(err)
		return []*models.MealToFront{}, internal.ErrReturningAllMeals
	}
	return
}

func (e *Endpoints) GetMeal(userId, mealId string) (meal models.MealToFront, err error) {
	request, err := http.NewRequest(http.MethodGet, config.Config.MealsURL+"user/"+userId+"/meal/"+mealId, nil)
	if err != nil {
		log.Error(err)
		return models.MealToFront{}, internal.ErrReturningMeal
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		newError := new(internal.ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return models.MealToFront{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&meal)
	if err != nil {
		log.Error(err)
		return models.MealToFront{}, internal.ErrReturningMeal
	}

	return
}

func getSeason() string {
	t := time.Now()
	switch t.Month() {
	case time.January, time.February, time.March:
		return "invierno"
	case time.April, time.May, time.June:
		return "primavera"
	case time.July, time.August, time.September:
		return "verano"
	case time.October, time.November, time.December:
		return "otoño"
	}
	return ""
}
