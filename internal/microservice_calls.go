package internal

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

const mealUrl = "http://localhost:3200/"
const userUrl = "http://localhost:3100/"

type Endpoints struct {
}
type EndpointsI interface {
	GetUser(userId string) (user User, err error)
	GetAllMeals(userId string) (meals []*MealToFront, err error)
	GetMeal(userId, mealId string) (meal MealToFront, err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetUser(userId string) (user User, err error) {
	request, err := http.NewRequest(http.MethodGet, userUrl+"user/"+userId, nil)
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
		log.Error(fmt.Sprintf("Error %d: %s", response.StatusCode, response.Status))
		return User{}, ErrReturningUser
	}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		log.Error(err)
		return User{}, ErrReturningUser
	}

	return
}

func (e *Endpoints) GetAllMeals(userId string) (meals []*MealToFront, err error) {

	request, err := http.NewRequest(http.MethodGet, mealUrl+"user/"+userId+"/meal", nil)
	if err != nil {
		log.Error(err)
		return []*MealToFront{}, ErrReturningAllMeals
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		log.Error(fmt.Sprintf("Error %d: %s", response.StatusCode, response.Status))
		return []*MealToFront{}, ErrReturningAllMeals
	}
	err = json.NewDecoder(response.Body).Decode(&meals)
	if err != nil {
		log.Error(err)
		return []*MealToFront{}, ErrReturningAllMeals
	}
	return
}

func (e *Endpoints) GetMeal(userId, mealId string) (meal MealToFront, err error) {
	request, err := http.NewRequest(http.MethodGet, mealUrl+"user/"+userId+"/meal/"+mealId, nil)
	if err != nil {
		log.Error(err)
		return MealToFront{}, ErrReturningMeal
	}
	response, err := httpClient.Do(request)
	if response.StatusCode > 299 {
		log.Error(fmt.Sprintf("Error %d: %s", response.StatusCode, response.Status))
		return MealToFront{}, ErrReturningMeal

	}
	err = json.NewDecoder(response.Body).Decode(&meal)
	if err != nil {
		log.Error(err)
		return MealToFront{}, ErrReturningMeal
	}

	return
}
