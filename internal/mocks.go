package internal

import (
	"github.com/stretchr/testify/mock"
)

type EndpointsMock struct {
	mock.Mock
}

func (e *EndpointsMock) GetUser(userId string) (user User, err error) {
	args := e.Called(userId)
	return args.Get(0).(User), args.Error(1)
}

func (e *EndpointsMock) GetAllMeals(userId string) (meals []*MealToFront, err error) {
	args := e.Called(userId)
	return args.Get(0).([]*MealToFront), args.Error(1)
}

func (e *EndpointsMock) GetMeal(userId, mealId string) (meal MealToFront, err error) {
	args := e.Called(userId, mealId)
	return args.Get(0).(MealToFront), args.Error(1)
}
