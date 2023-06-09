package internal

import (
	"calendar/pkg/database"
	"github.com/labstack/gommon/log"
)

const (
	getCalendar    = "SELECT * FROM calendar WHERE user_id = ? ORDER BY date"
	updateCalendar = "UPDATE calendar SET meal_id = ? WHERE user_id = ? AND date = ?"
	createCalendar = "INSERT INTO calendar (meal_id,user_id,date) VALUES (?,?,?)"
	deleteCalendar = "DELETE FROM calendar WHERE user_id = ?"

	specificDateCalendar = "SELECT * FROM calendar WHERE user_id = ? AND date = ?"
)

type SQLiteCalendarRepository struct {
	db *database.Database
}

type DBCalendarI interface {
	GetCalendar(id string) (calendar []Calendar, err error)
	UpdateCalendar(id string, calendar CalendarUpdate) (err error)
	CreateCalendar(calendar []Calendar) (err error)
	DeleteCalendar(id string) (err error)

	GetCalendarSpecificDate(id, date string) (calendar []Calendar, err error)
}

func NewSQLiteCalendarRepository(db *database.Database) *SQLiteCalendarRepository {
	return &SQLiteCalendarRepository{
		db: db,
	}
}

func (r *SQLiteCalendarRepository) GetCalendar(id string) (calendar []Calendar, err error) {
	err = r.db.Conn.Select(&calendar, getCalendar, id)
	if err != nil {
		log.Error(err)
		return
	}
	if len(calendar) == 0 {
		err = ErrCalendarNotFound
	}
	return
}

func (r *SQLiteCalendarRepository) UpdateCalendar(id string, c CalendarUpdate) (err error) {
	_, err = r.db.Conn.Exec(updateCalendar, c.MealId, id, c.MealDate)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (r *SQLiteCalendarRepository) CreateCalendar(calendar []Calendar) (err error) {
	for _, c := range calendar {
		_, err = r.db.Conn.Exec(createCalendar, c.MealId, c.UserId, c.Date)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

func (r *SQLiteCalendarRepository) DeleteCalendar(id string) (err error) {
	_, err = r.db.Conn.Exec(deleteCalendar, id)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (r *SQLiteCalendarRepository) GetCalendarSpecificDate(id, date string) (calendar []Calendar, err error) {
	err = r.db.Conn.Select(&calendar, specificDateCalendar, id, date)
	if err != nil {
		log.Error(err)
		return
	}
	if len(calendar) == 0 {
		err = ErrDateNotFound
	}
	return
}
