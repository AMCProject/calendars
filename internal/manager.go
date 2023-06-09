package internal

import (
	"calendar/pkg/database"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type ICalendarManager interface {
	GetCalendar(id string) (calendar []Calendar, err error)
	UpdateCalendar(id string, calendar []CalendarUpdate) (calendarResponse []Calendar, err error)
	CreateCalendar(id string) (calendar []Calendar, err error)
	DeleteCalendar(id string) (err error)
	GetFrontCalendar(calendar []Calendar) (finalCal []CalendarWithMeals, err error)
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
	if _, err = Microservices.GetUser(id); err != nil {
		return
	}
	calendar, err = c.db.GetCalendar(id)
	if err != nil {
		return
	}
	if len(calendar) > 0 {
		t := time.Now()
		t = t.AddDate(0, 0, int(21+(7-t.Weekday())))
		tFormat := t.Format("2006-01-02")
		if !strings.EqualFold(tFormat, calendar[len(calendar)-1].Date) {
			meals, errM := Microservices.GetAllMeals(id)
			if errM != nil {
				return calendar, errM
			}
			lastD, errF := time.Parse("2006-01-02", calendar[len(calendar)-1].Date)
			if errF != nil {
				return calendar, errF
			}

			days := int(t.Sub(lastD).Hours() / 24)
			if days > 28 {
				calendar, err = c.utils.CalendarCreator(id, meals)
			}
			//_ = s.Repository.DeleteCalendar(id)
			calendar, err = c.utils.UpdateNewDays(id, calendar, meals, days)
		}
		_ = c.db.DeleteCalendar(id)
		_ = c.db.CreateCalendar(calendar)
	}
	return
}

func (c *CalendarManager) UpdateCalendar(id string, calendar []CalendarUpdate) (calendarResponse []Calendar, err error) {
	if _, err = Microservices.GetUser(id); err != nil {
		return
	}
	for _, cal := range calendar {
		if _, err = c.db.GetCalendarSpecificDate(id, cal.MealDate); err != nil {
			return
		}
	}
	for _, cal := range calendar {
		if err = c.db.UpdateCalendar(id, cal); err != nil {
			return
		}
	}

	return c.db.GetCalendar(id)
}

func (c *CalendarManager) CreateCalendar(id string) (calendar []Calendar, err error) {
	if _, err = c.db.GetCalendar(id); err == nil {
		return []Calendar{}, ErrCalendarAlreadyExists
	}
	if _, err = Microservices.GetUser(id); err != nil {
		return
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
	if _, err = Microservices.GetUser(id); err != nil {
		return
	}
	if _, err = c.db.GetCalendar(id); err != nil {
		return
	}
	return c.db.DeleteCalendar(id)
}

func (c *CalendarManager) GetFrontCalendar(calendar []Calendar) (finalCal []CalendarWithMeals, err error) {
	diff := 28 - len(calendar)
	for i := 0; i < diff; i++ {
		calAux := CalendarWithMeals{MealId: "", Name: "NO HAY COMIDA"}
		finalCal = append(finalCal, calAux)
	}
	for _, cal := range calendar {
		var meal MealToFront
		meal, err = Microservices.GetMeal(cal.UserId, cal.MealId)
		if err != nil {
			return
		}
		calAux := CalendarWithMeals{MealId: cal.MealId, UserId: cal.UserId, Date: cal.Date, Name: meal.Name}
		finalCal = append(finalCal, calAux)
	}
	return
}
