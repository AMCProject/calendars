package utils

import (
	"calendar/internal/models"
	"math"
	"math/rand"
	"strings"
	"time"
)

type CalendarTools struct{}

type ICalendarTools interface {
	CalendarCreator(userId string, meals []*models.MealToFront) (calendar []models.Calendar, err error)
	UpdateDaysInCalendar(d string, calendar []models.Calendar, meals []*models.MealToFront, dates models.UpdateWeekCalendar) (finalCalendar []models.Calendar, err error)
	UpdateNewDays(userId string, calendar []models.Calendar, meals []*models.MealToFront, days int) (finalCalendar []models.Calendar, err error)
	ReturnRandomMeal(calendar []models.Calendar, meals []*models.MealToFront, wd int) (meal models.MealToFront)
	CalendarContains(calendar []models.Calendar, mealId string) (distance float64)
	SpecialMeal(meal *models.MealToFront, numb float64, wd int) (res float64)
	GetHighestMeal(keyMeal []float64) (index int)
}

func NewCalendarToolsManager() *CalendarTools {
	return &CalendarTools{}
}

func (s *CalendarTools) CalendarCreator(userId string, meals []*models.MealToFront) (calendar []models.Calendar, err error) {
	var days int
	t := time.Now()
	wd := t.Weekday()
	if wd == 0 {
		days = 21
	} else {
		days = 21 + (7 - int(wd))
	}
	for i := 0; i <= days; i++ {
		newDate := t.AddDate(0, 0, i)
		meal := s.ReturnRandomMeal(calendar, meals, newDate)
		cal := models.Calendar{
			UserId: userId,
			MealId: meal.Id,
			Name:   meal.Name,
			Date:   newDate.Format("2006/01/02"),
		}
		calendar = append(calendar, cal)
	}

	return
}

func (s *CalendarTools) UpdateDaysInCalendar(id string, calendar []models.Calendar, meals []*models.MealToFront, dates models.UpdateWeekCalendar) (finalCalendar []models.Calendar, err error) {
	var inRange bool
	finalCalendar = calendar
	for i, c := range finalCalendar {
		if c.Date == dates.From {
			inRange = true
		}
		if !inRange {
			continue
		}
		updateDay, _ := time.Parse("2006/01/02", c.Date)
		meal := s.ReturnRandomMeal(finalCalendar, meals, updateDay)
		finalCalendar[i] = models.Calendar{
			UserId: id,
			MealId: meal.Id,
			Name:   meal.Name,
			Date:   c.Date,
		}
		if c.Date == dates.To {
			inRange = false
		}
	}
	return
}

func (s *CalendarTools) UpdateNewDays(userId string, calendar []models.Calendar, meals []*models.MealToFront, days int) (finalCalendar []models.Calendar, err error) {
	finalCalendar = calendar
	if len(calendar) >= 28 {
		finalCalendar = calendar[days:]
	}
	t, _ := time.Parse("2006/01/02", finalCalendar[len(finalCalendar)-1].Date)
	for i := 0; i < days; i++ {
		newDate := t.AddDate(0, 0, i+1)
		meal := s.ReturnRandomMeal(finalCalendar, meals, newDate)
		cal := models.Calendar{
			UserId: userId,
			MealId: meal.Id,
			Name:   meal.Name,
			Date:   newDate.Format("2006/01/02"),
		}
		finalCalendar = append(finalCalendar, cal)

	}
	return
}

func (s *CalendarTools) ReturnRandomMeal(calendar []models.Calendar, meals []*models.MealToFront, date time.Time) (meal models.MealToFront) {
	var keyMeal []float64
	for _, m := range meals {
		numb := math.Abs(rand.Float64() * 3)
		contains, distance := s.CalendarContains(calendar, m.Id, date)
		if distance == 1 {
			numb = numb - 20
		}
		if distance > 0 {
			numb = numb - 1.9 + ((distance / float64(len(calendar))) / 4)
		}
		if distance == 0 && !contains {
			numb += 0.8
		}
		if distance == 0 && contains {
			numb = numb - 20
		}
		numb = s.SpecialMeal(m, numb, int(date.Weekday()))
		if strings.EqualFold(m.Type, models.Semanal) && (distance >= 7 || distance == 0) {
			if distance == 0 {
				numb += 1.2
			} else {
				numb += 1.6 - (1/(distance))*2.3
			}
		}
		keyMeal = append(keyMeal, numb)
	}
	index := s.GetHighestMeal(keyMeal)
	meal = *meals[index]
	return
}

func (s *CalendarTools) CalendarContains(calendar []models.Calendar, mealId string, date time.Time) (contains bool, distance float64) {
	distance = 100
	for _, c := range calendar {
		if c.MealId == mealId {
			contains = true
			compareDate, _ := time.Parse("2006/01/02", c.Date)
			difference := math.Abs(date.Sub(compareDate).Hours() / 24)
			if difference < distance {
				distance = difference
			}
		}
	}
	if !contains {
		distance = 0
	}
	return
}

func (s *CalendarTools) SpecialMeal(meal *models.MealToFront, numb float64, wd int) (res float64) {
	res = numb
	if strings.EqualFold(meal.Type, models.Ocasional) && (wd == 0 || wd == 6) {
		res += 2.10
	}
	if strings.EqualFold(meal.Type, models.Ocasional) && (wd > 0 && wd < 6) {
		res -= 2.9
	}
	return res
}

func (s *CalendarTools) GetHighestMeal(meals []float64) (index int) {
	highest := meals[0]
	for i, m := range meals {
		if m > highest {
			highest = m
			index = i
		}
	}
	return
}
