package managers

import (
	"calendar/internal"
	"calendar/internal/models"
	"calendar/internal/repositories"
	"calendar/internal/utils"
	"calendar/pkg/database"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"
)

type ICalendarManager interface {
	GetCalendar(id string) (calendar []models.Calendar, err error)
	UpdateCalendar(id string, calendar models.Calendar) (calendarResponse []models.Calendar, err error)
	UpdateDaysCalendar(id string, dates models.UpdateWeekCalendar) (calendar []models.Calendar, err error)
	CreateCalendar(id string) (calendar []models.Calendar, err error)
	DeleteCalendar(id string) (err error)
	GetFrontCalendar(calendar []models.Calendar) (finalCal []models.Calendar, err error)
}

var Microservices utils.EndpointsI = &utils.Endpoints{}

type CalendarManager struct {
	db       *repositories.SQLiteCalendarRepository
	validate *validator.Validate
	utils    *utils.CalendarTools
}

func NewCalendarManager(db database.Database) *CalendarManager {
	return &CalendarManager{
		db:       repositories.NewSQLiteCalendarRepository(&db),
		validate: validator.New(),
		utils:    utils.NewCalendarToolsManager(),
	}
}

func (c *CalendarManager) GetCalendar(id string) (calendar []models.Calendar, err error) {
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
		if len(meals) == 0 {
			if err = c.db.DeleteCalendar(id); err != nil {
				return []models.Calendar{}, internal.ErrSomethingWentWrong
			}
			return []models.Calendar{}, internal.ErrMealsNotFound
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
			calendar = calendar[len(calendar)-28:]
		}
		if err = c.db.DeleteCalendar(id); err != nil {
			return []models.Calendar{}, internal.ErrSomethingWentWrong
		}
		if err = c.db.CreateCalendar(calendar); err != nil {
			return []models.Calendar{}, internal.ErrSomethingWentWrong
		}
	}

	return
}

func (c *CalendarManager) UpdateCalendar(id string, calendar models.Calendar) (calendarResponse []models.Calendar, err error) {
	var meal models.MealToFront
	_, err = time.Parse("2006/01/02", calendar.Date)
	if err != nil {
		return []models.Calendar{}, internal.ErrInvalidDateFormat
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

func (c *CalendarManager) UpdateDaysCalendar(id string, dates models.UpdateWeekCalendar) (calendar []models.Calendar, err error) {
	calendar, err = c.db.GetCalendar(id)
	if err != nil {
		return []models.Calendar{}, err
	}
	_, err = time.Parse("2006/01/02", dates.From)
	if err != nil {
		return []models.Calendar{}, internal.ErrInvalidDateFormat
	}
	_, err = time.Parse("2006/01/02", dates.To)
	if err != nil {
		return []models.Calendar{}, internal.ErrInvalidDateFormat
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
	if len(meals) == 0 {
		if err = c.db.DeleteCalendar(id); err != nil {
			return []models.Calendar{}, internal.ErrSomethingWentWrong
		}
		return []models.Calendar{}, internal.ErrMealsNotFound
	}
	finalCal, err := c.utils.UpdateDaysInCalendar(id, calendar, meals, dates)
	if err != nil {
		return []models.Calendar{}, err
	}
	if err = c.db.DeleteCalendar(id); err != nil {
		return []models.Calendar{}, internal.ErrSomethingWentWrong
	}
	if err = c.db.CreateCalendar(finalCal); err != nil {
		return []models.Calendar{}, internal.ErrSomethingWentWrong
	}
	return finalCal, err
}

func (c *CalendarManager) CreateCalendar(id string) (calendar []models.Calendar, err error) {
	if _, err = c.db.GetCalendar(id); err == nil {
		return []models.Calendar{}, internal.ErrCalendarAlreadyExists
	}
	meals, err := Microservices.GetAllMeals(id)
	if len(meals) == 0 {
		if err = c.db.DeleteCalendar(id); err != nil {
			return []models.Calendar{}, internal.ErrSomethingWentWrong
		}
		return []models.Calendar{}, internal.ErrMealsNotFound
	}
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

func (c *CalendarManager) GetFrontCalendar(calendar []models.Calendar) (finalCal []models.Calendar, err error) {
	diff := 28 - len(calendar)
	firstDate, _ := time.Parse("2006/01/02", calendar[0].Date)
	for i := 0; i < diff; i++ {
		noMealDate := firstDate.AddDate(0, 0, -(diff - i))
		calAux := models.Calendar{MealId: "", Name: "NO MEAL", Date: noMealDate.Format("2006/01/02")}
		finalCal = append(finalCal, calAux)
	}
	for _, cal := range calendar {
		calAux := models.Calendar{MealId: cal.MealId, UserId: cal.UserId, Date: cal.Date, Name: cal.Name}
		finalCal = append(finalCal, calAux)
	}
	return
}
