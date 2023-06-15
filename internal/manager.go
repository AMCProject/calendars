package internal

import (
	"calendar/pkg/database"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type ICalendarManager interface {
	GetCalendar(id string) (calendar []Calendar, err error)
	UpdateCalendar(id string, calendar Calendar) (calendarResponse []Calendar, err error)
	UpdateDaysCalendar(id string, dates UpdateWeekCalendar) (calendar []Calendar, err error)
	CreateCalendar(id string) (calendar []Calendar, err error)
	DeleteCalendar(id string) (err error)
	GetFrontCalendar(calendar []Calendar) (finalCal []Calendar, err error)
}

var Microservices EndpointsI = &Endpoints{}

type CalendarManager struct {
	db       *SQLiteCalendarRepository
	validate *validator.Validate
	utils    *CalendarTools
}

func NewCalendarManager(db database.Database) *CalendarManager {
	return &CalendarManager{
		db:       NewSQLiteCalendarRepository(&db),
		validate: validator.New(),
		utils:    NewCalendarToolsManager(),
	}
}

func (c *CalendarManager) GetCalendar(id string) (calendar []Calendar, err error) {
	calendar, err = c.db.GetCalendar(id)
	if err != nil {
		return
	}

	t := time.Now()
	wd := t.Weekday()
	var differenceDays int
	if wd == 0 {
		differenceDays = 21
	} else {
		differenceDays = 21 + (7 - int(wd))
	}
	t = t.AddDate(0, 0, differenceDays)
	tFormat := t.Format("2006/01/02")
	if !strings.EqualFold(tFormat, calendar[len(calendar)-1].Date) {
		meals, errM := Microservices.GetAllMeals(id)
		if errM != nil {
			return calendar, errM
		}
		lastD, errF := time.Parse("2006/01/02", calendar[len(calendar)-1].Date)
		if errF != nil {
			return calendar, errF
		}
		days := int(t.Sub(lastD).Hours() / 24)
		if days > 28 {
			calendar, err = c.utils.CalendarCreator(id, meals)
		}
		//_ = s.Repository.DeleteCalendar(id)
		calendar, err = c.utils.UpdateNewDays(id, calendar, meals, days)
		if len(calendar) > 28 {
			calendar = calendar[len(calendar)+28:]
		}
		if err = c.db.DeleteCalendar(id); err != nil {
			return []Calendar{}, ErrSomethingWentWrong
		}
		if err = c.db.CreateCalendar(calendar); err != nil {
			return []Calendar{}, ErrSomethingWentWrong
		}
	}

	return
}

func (c *CalendarManager) UpdateCalendar(id string, calendar Calendar) (calendarResponse []Calendar, err error) {
	var meal MealToFront
	_, err = time.Parse("2006/01/02", calendar.Date)
	if err != nil {
		return []Calendar{}, ErrInvalidDateFormat
	}
	if _, err = c.db.GetCalendarSpecificDate(id, calendar.Date); err != nil {
		return
	}

	if meal, err = Microservices.GetMeal(id, calendar.MealId); err != nil {
		return
	}
	calendar.Name = meal.Name

	if err = c.db.UpdateCalendar(id, calendar); err != nil {
		return
	}

	return c.db.GetCalendar(id)
}

func (c *CalendarManager) UpdateDaysCalendar(id string, dates UpdateWeekCalendar) (calendar []Calendar, err error) {
	calendar, err = c.db.GetCalendar(id)
	if err != nil {
		return []Calendar{}, err
	}
	_, err = time.Parse("2006/01/02", dates.From)
	if err != nil {
		return []Calendar{}, ErrInvalidDateFormat
	}
	_, err = time.Parse("2006/01/02", dates.To)
	if err != nil {
		return []Calendar{}, ErrInvalidDateFormat
	}
	if _, err = c.db.GetCalendarSpecificDate(id, dates.From); err != nil {
		return nil, err
	}
	if _, err = c.db.GetCalendarSpecificDate(id, dates.To); err != nil {
		return nil, err
	}
	meals, err := Microservices.GetAllMeals(id)
	if err != nil {
		return
	}
	finalCal, err := c.utils.UpdateDaysInCalendar(id, calendar, meals, dates)
	if err != nil {
		return []Calendar{}, err
	}
	if err = c.db.DeleteCalendar(id); err != nil {
		return []Calendar{}, ErrSomethingWentWrong
	}
	if err = c.db.CreateCalendar(finalCal); err != nil {
		return []Calendar{}, ErrSomethingWentWrong
	}
	return finalCal, err
}

func (c *CalendarManager) CreateCalendar(id string) (calendar []Calendar, err error) {
	if _, err = c.db.GetCalendar(id); err == nil {
		return []Calendar{}, ErrCalendarAlreadyExists
	}
	meals, err := Microservices.GetAllMeals(id)
	if err != nil {
		return
	}
	calendar, err = c.utils.CalendarCreator(id, meals)
	if err != nil {
		return
	}
	err = c.db.CreateCalendar(calendar)
	if err != nil {
		return
	}
	return
}

func (c *CalendarManager) DeleteCalendar(id string) (err error) {
	if _, err = c.db.GetCalendar(id); err != nil {
		return
	}
	return c.db.DeleteCalendar(id)
}

func (c *CalendarManager) GetFrontCalendar(calendar []Calendar) (finalCal []Calendar, err error) {
	diff := 28 - len(calendar)
	firstDate, _ := time.Parse("2006/01/02", calendar[0].Date)
	for i := 0; i < diff; i++ {
		noMealDate := firstDate.AddDate(0, 0, -(diff - i))
		calAux := Calendar{MealId: "", Name: "NO MEAL", Date: noMealDate.Format("2006/01/02")}
		finalCal = append(finalCal, calAux)
	}
	for _, cal := range calendar {
		calAux := Calendar{MealId: cal.MealId, UserId: cal.UserId, Date: cal.Date, Name: cal.Name}
		finalCal = append(finalCal, calAux)
	}
	return
}
