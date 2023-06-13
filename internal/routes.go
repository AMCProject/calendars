package internal

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	RouteCalendar         = "/user/:user_id/calendar"
	RouteCalendarRedo     = "/user/:user_id/redo"
	RouteCalendarRedoWeek = "/user/:user_id/redoweek"

	ParamUserID = "user_id"
)

type ErrorResponse struct {
	Err ErrorBody `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Err.Status, e.Err.Message)
}

type ErrorBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorResponse(c echo.Context, err error) error {
	errResponse := &ErrorResponse{Err: errorsMap[err.Error()]}
	if errResponse.Err.Status == 0 {
		if err := c.JSON(http.StatusInternalServerError, err); err != nil {
			return err
		}
		return err
	}
	if err := c.JSON(errResponse.Err.Status, errResponse); err != nil {
		return err
	}
	return errResponse
}

var errorsMap = map[string]ErrorBody{
	ErrUserIDNotPresent.Error():      {Status: http.StatusBadRequest, Message: ErrUserIDNotPresent.Error()},
	ErrWrongBody.Error():             {Status: http.StatusBadRequest, Message: ErrWrongBody.Error()},
	ErrInvalidDateFormat.Error():     {Status: http.StatusBadRequest, Message: ErrInvalidDateFormat.Error()},
	ErrCalendarNotFound.Error():      {Status: http.StatusNotFound, Message: ErrCalendarNotFound.Error()},
	ErrUserNotFound.Error():          {Status: http.StatusNotFound, Message: ErrUserNotFound.Error()},
	ErrMealNotFound.Error():          {Status: http.StatusNotFound, Message: ErrMealNotFound.Error()},
	ErrDateNotFound.Error():          {Status: http.StatusNotFound, Message: ErrDateNotFound.Error()},
	ErrCalendarAlreadyExists.Error(): {Status: http.StatusConflict, Message: ErrCalendarAlreadyExists.Error()},
	ErrSomethingWentWrong.Error():    {Status: http.StatusInternalServerError, Message: ErrSomethingWentWrong.Error()},
	ErrReturningAllMeals.Error():     {Status: http.StatusInternalServerError, Message: ErrReturningAllMeals.Error()},
	ErrReturningMeal.Error():         {Status: http.StatusInternalServerError, Message: ErrReturningMeal.Error()},
	ErrReturningUser.Error():         {Status: http.StatusInternalServerError, Message: ErrReturningUser.Error()},
}
var (
	ErrUserIDNotPresent      = errors.New("error with userID given")
	ErrSomethingWentWrong    = errors.New("something went wrong")
	ErrWrongBody             = errors.New("malformed body")
	ErrCalendarNotFound      = errors.New("calendar not found")
	ErrCalendarAlreadyExists = errors.New("this user already has a calendar")
	ErrUserNotFound          = errors.New("user not found")
	ErrMealNotFound          = errors.New("meal not found")
	ErrReturningAllMeals     = errors.New("unexpect error recovering meals")
	ErrReturningMeal         = errors.New("unexpect error recovering specific meal")
	ErrReturningUser         = errors.New("unexpect error recovering user")
	ErrDateNotFound          = errors.New("date not found in this calendar")
	ErrInvalidDateFormat     = errors.New("invalid date format, must be dd/MM/yyyy")
)
