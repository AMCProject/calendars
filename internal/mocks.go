package internal

import (
	"calendar/internal/models"
	"github.com/stretchr/testify/mock"
)

type EndpointsMock struct {
	mock.Mock
}

func (e *EndpointsMock) GetAllMeals(userId string) (meals []*models.MealToFront, err error) {
	args := e.Called(userId)
	return args.Get(0).([]*models.MealToFront), args.Error(1)
}

func (e *EndpointsMock) GetMeal(userId, mealId string) (meal models.MealToFront, err error) {
	args := e.Called(userId, mealId)
	return args.Get(0).(models.MealToFront), args.Error(1)
}
