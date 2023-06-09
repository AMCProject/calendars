package internal

import (
	"calendar/pkg/database"
	"calendar/pkg/url"

	"github.com/labstack/echo/v4"

	"net/http"
)

type CalendarAPI struct {
	DB      database.Database
	Manager ICalendarManager
	Utils   ICalendarTools
}

func (a *CalendarAPI) PostCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}
	var calendarOutput []CalendarOutput
	calendar, err := a.Manager.CreateCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	for _, cal := range calendar {
		meal, _ := Microservices.GetMeal(cal.UserId, cal.MealId)
		co := CalendarOutput{Name: meal.Name, Date: cal.Date}
		calendarOutput = append(calendarOutput, co)
	}

	return c.JSON(http.StatusCreated, calendarOutput)
}

func (a *CalendarAPI) GetCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}
	calendar, err := a.Manager.GetCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)
}

func (a *CalendarAPI) PutCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	calendarReq := &CalendarUpdate{}
	if err := c.Bind(calendarReq); err != nil {
		return NewErrorResponse(c, ErrWrongBody)
	}

	calendar, err := a.Manager.UpdateCalendar(userID, *calendarReq)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, calendar)
}

func (a *CalendarAPI) DeleteCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}
	err := a.Manager.DeleteCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.NoContent(http.StatusNoContent)

}

func (a *CalendarAPI) RedoCalendarHandler(c echo.Context) error {

	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	err := a.Manager.DeleteCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	calendar, err := a.Manager.CreateCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	var calendarOutput []CalendarOutput

	for _, cal := range calendar {
		meal, _ := Microservices.GetMeal(cal.UserId, cal.MealId)
		co := CalendarOutput{Name: meal.Name, Date: cal.Date}
		calendarOutput = append(calendarOutput, co)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)

}
