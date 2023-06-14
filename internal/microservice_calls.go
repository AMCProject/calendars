package internal

import (
	"calendar/internal/config"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type Endpoints struct {
}
type EndpointsI interface {
	GetUser(userId string) (user User, err error)
	GetAllMeals(userId string) (meals []*MealToFront, err error)
	GetMeal(userId, mealId string) (meal MealToFront, err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetUser(userId string) (user User, err error) {
	request, err := http.NewRequest(http.MethodGet, config.Config.UsersURL+"user/"+userId, nil)
	if err != nil {
		log.Error(err)
		return User{}, ErrReturningUser
	}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error(err)
		return User{}, ErrReturningUser
	}
	if response.StatusCode > 299 {
		newError := new(ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return User{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		log.Error(err)
		return User{}, ErrReturningUser
	}

	return
}

func (e *Endpoints) GetAllMeals(userId string) (meals []*MealToFront, err error) {
	url := config.Config.MealsURL + "user/" + userId + "/meal"
	season := getSeason()
	if season != "" {
		url += "?season[]=" + season
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return []*MealToFront{}, ErrReturningAllMeals
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		newError := new(ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return []*MealToFront{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&meals)
	if err != nil {
		log.Error(err)
		return []*MealToFront{}, ErrReturningAllMeals
	}
	return
}

func (e *Endpoints) GetMeal(userId, mealId string) (meal MealToFront, err error) {
	request, err := http.NewRequest(http.MethodGet, config.Config.MealsURL+"user/"+userId+"/meal/"+mealId, nil)
	if err != nil {
		log.Error(err)
		return MealToFront{}, ErrReturningMeal
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		newError := new(ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return MealToFront{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&meal)
	if err != nil {
		log.Error(err)
		return MealToFront{}, ErrReturningMeal
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
