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
	ErrUserIDNotPresent      = errors.New("error con el ID del usuario dado")
	ErrSomethingWentWrong    = errors.New("error inesperado")
	ErrWrongBody             = errors.New("el cuerpo enviado es erróneo")
	ErrCalendarNotFound      = errors.New("calendario no encontrado")
	ErrCalendarAlreadyExists = errors.New("este usuario ya tiene un calendario")
	ErrUserNotFound          = errors.New("usuario no encontrado")
	ErrMealNotFound          = errors.New("comida no encontrada")
	ErrReturningAllMeals     = errors.New("error inesperado recuperando las comidas")
	ErrReturningMeal         = errors.New("error inesperado recuperando la información de la comida")
	ErrReturningUser         = errors.New("error inesperado recuperando la información del usuario")
	ErrDateNotFound          = errors.New("fecha indicada no encontrada en el calendario")
	ErrInvalidDateFormat     = errors.New("formato inválido de fecha, debe ser aaaa/MM/dd")
)
