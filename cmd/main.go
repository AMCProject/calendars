package main

import (
	"calendar/internal"
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/managers"
	"calendar/pkg/database"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

const (
	banner = `
   ___    __  ___  _____       _____         __
  / _ |  /  |/  / / ___/      / ___/ ___ _  / / ___   ___  ___ _  ____  ___
 / __ | / /|_/ / / /__       / /__  / _  / / / / -_) / _ \/ _  / / __/ (_-<
/_/ |_|/_/  /_/  \___/       \___/  \_,_/ /_/  \__/ /_//_/\_,_/ /_/   /___/

AMC Calendars Service
`
)

func main() {
	if err := config.LoadConfiguration(); err != nil {
		log.Fatal(err)
	}
	db := database.InitDB(config.Config.DBName)
	e := setUpServer(db)
	e.Logger.Fatal(e.Start(config.Config.Host + ":" + config.Config.Port))

}

func setUpServer(db *database.Database) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	addRoutes(e, *db)
	e.HideBanner = true
	fmt.Printf(banner)

	return e

}

func addRoutes(e *echo.Echo, db database.Database) {

	calendarManager := managers.NewCalendarManager(db)

	calendarAPI := handlers.CalendarAPI{DB: db, Manager: calendarManager}
	e.GET(internal.RouteCalendar, calendarAPI.GetCalendarHandler)
	e.POST(internal.RouteCalendar, calendarAPI.PostCalendarHandler)
	e.PUT(internal.RouteCalendar, calendarAPI.PutCalendarHandler)
	e.DELETE(internal.RouteCalendar, calendarAPI.DeleteCalendarHandler)

	e.PUT(internal.RouteCalendarRedo, calendarAPI.RedoCalendarHandler)
	e.PUT(internal.RouteCalendarRedoWeek, calendarAPI.RedoWeekCalendarHandler)
}
