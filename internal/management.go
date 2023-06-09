package internal

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

type CalendarTools struct{}

type ICalendarTools interface {
	CalendarCreator(userId string, meals []*MealToFront) (calendar []Calendar, err error)
	UpdateNewDays(userId string, calendar []Calendar, meals []*MealToFront, days int) (finalCalendar []Calendar, err error)
	ReturnRandomMeal(calendar []Calendar, meals []*MealToFront, wd int) (meal MealToFront)
	CalendarContains(calendar []Calendar, mealId string) (distance float64)
	SpecialMeal(meal *MealToFront, numb float64, wd int) (res float64)
	GetHighestMeal(keyMeal []float64) (index int)
}

func NewCalendarToolsManager() *CalendarTools {
	return &CalendarTools{}
}

func (s *CalendarTools) CalendarCreator(userId string, meals []*MealToFront) (calendar []Calendar, err error) {
	var days int
	t := time.Now()
	if t.Hour() >= 14 {
		t = t.AddDate(0, 0, 1)
	}
	wd := t.Weekday()
	if wd == 0 {
		days = 28
	} else {
		days = 21 + (7 - int(wd))
	}
	for i := 0; i <= days; i++ {
		newDate := t.AddDate(0, 0, i)
		meal := s.ReturnRandomMeal(calendar, meals, int(newDate.Weekday()))
		cal := Calendar{
			UserId:   userId,
			MealId:   meal.Id,
			MealName: meal.Name,
			Date:     newDate.Format("2006-01-02"),
		}
		calendar = append(calendar, cal)
	}

	return
}

func (s *CalendarTools) UpdateNewDays(userId string, calendar []Calendar, meals []*MealToFront, days int) (finalCalendar []Calendar, err error) {
	finalCalendar = calendar
	if len(calendar) >= 28 {
		finalCalendar = calendar[days:]
	}
	t, _ := time.Parse("2006-01-02", finalCalendar[len(finalCalendar)-1].Date)
	for i := 0; i < days; i++ {
		newDate := t.AddDate(0, 0, i+1)
		meal := s.ReturnRandomMeal(finalCalendar, meals, int(newDate.Weekday()))
		cal := Calendar{
			UserId: userId,
			MealId: meal.Id,
			Date:   newDate.Format("2006-01-02"),
		}
		finalCalendar = append(finalCalendar, cal)

	}
	return
}

func (s *CalendarTools) ReturnRandomMeal(calendar []Calendar, meals []*MealToFront, wd int) (meal MealToFront) {
	var keyMeal []float64
	for _, m := range meals {
		numb := math.Abs(rand.Float64() * 3)
		distance := s.CalendarContains(calendar, m.Id)
		if distance > 0 {
			numb = numb - 1.9 + ((distance / float64(len(calendar))) / 4)
		}
		if distance == 0 {
			numb += 0.8
		}
		numb = s.SpecialMeal(m, numb, wd)
		if strings.EqualFold(m.Type, Semanal) && (distance >= 7 || distance == 0) {
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

func (s *CalendarTools) CalendarContains(calendar []Calendar, mealId string) (distance float64) {
	var contains bool
	for i, c := range calendar {
		if c.MealId == mealId {
			contains = true
			distance = float64(i)
		}
	}
	if contains {
		distance = float64(len(calendar)) - distance
	}
	return
}

func (s *CalendarTools) SpecialMeal(meal *MealToFront, numb float64, wd int) (res float64) {
	res = numb
	if strings.EqualFold(meal.Type, Ocasional) && (wd == 0 || wd == 6) {
		res += 2.10
	}
	if strings.EqualFold(meal.Type, Ocasional) && (wd > 0 && wd < 6) {
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
