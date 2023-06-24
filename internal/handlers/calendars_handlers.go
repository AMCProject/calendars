package handlers

import (
	"calendar/internal"
	"calendar/internal/managers"
	"calendar/internal/models"
	"calendar/internal/utils"
	"calendar/pkg/database"
	"calendar/pkg/url"

	"github.com/labstack/echo/v4"

	"net/http"
)

type CalendarAPI struct {
	DB      database.Database
	Manager managers.ICalendarManager
	Utils   utils.ICalendarTools
}

func (a *CalendarAPI) PostCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}
	calendar, err := a.Manager.CreateCalendar(userID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, finalCal)
}

func (a *CalendarAPI) GetCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}
	calendar, err := a.Manager.GetCalendar(userID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)
}

func (a *CalendarAPI) PutCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	calendarReq := &models.Calendar{}
	if err := c.Bind(calendarReq); err != nil {
		return internal.NewErrorResponse(c, internal.ErrWrongBody)
	}

	calendar, err := a.Manager.UpdateCalendar(userID, *calendarReq)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)
}

func (a *CalendarAPI) DeleteCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}
	err := a.Manager.DeleteCalendar(userID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	return c.NoContent(http.StatusNoContent)

}

func (a *CalendarAPI) RedoCalendarHandler(c echo.Context) error {

	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	err := a.Manager.DeleteCalendar(userID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	calendar, err := a.Manager.CreateCalendar(userID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}

	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, finalCal)

}

func (a *CalendarAPI) RedoWeekCalendarHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}
	dates := &models.UpdateWeekCalendar{}
	if err := c.Bind(dates); err != nil {
		return internal.NewErrorResponse(c, internal.ErrWrongBody)
	}

	calendar, err := a.Manager.UpdateDaysCalendar(userID, *dates)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}

	finalCal, err := a.Manager.GetFrontCalendar(calendar)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, finalCal)
}
