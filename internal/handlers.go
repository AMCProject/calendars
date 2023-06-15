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
	calendar, err := a.Manager.CreateCalendar(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, finalCal)
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

	calendarReq := &Calendar{}
	if err := c.Bind(calendarReq); err != nil {
		return NewErrorResponse(c, ErrWrongBody)
	}

	calendar, err := a.Manager.UpdateCalendar(userID, *calendarReq)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)
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

	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)

}

func (a *CalendarAPI) RedoWeekCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}
	dates := &UpdateWeekCalendar{}
	if err := c.Bind(dates); err != nil {
		return NewErrorResponse(c, ErrWrongBody)
	}

	calendar, err := a.Manager.UpdateDaysCalendar(userID, *dates)
	if err != nil {
		return NewErrorResponse(c, err)
	}

	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return NewErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, finalCal)
}
