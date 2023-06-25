package handlers

import (
	"bytes"
	"calendar/internal"
	"calendar/internal/managers"
	"calendar/internal/models"
	"calendar/pkg/database"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var databaseTest = "/amc_test.db"
var mealsDb = []*models.MealToFront{
	{
		Id:     "01FN3EEB2NVFJAHAPM00000001",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal1",
		Type:   "occasional",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000002",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal2",
		Type:   "semanal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000003",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal3",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000004",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal4",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000005",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal5",
		Type:   "occasional",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000006",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal6",
		Type:   "semanal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000007",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal7",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000008",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal8",
		Type:   "occasional",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000009",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal9",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000010",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal10",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000011",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal11",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000012",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal12",
		Type:   "normal",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000013",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal13",
		Type:   "occasional",
	},
	{
		Id:     "01FN3EEB2NVFJAHAPM00000014",
		UserId: "01FN3EEB2NVFJAHAPU00000001",
		Name:   "meal14",
		Type:   "normal",
	},
}

type CalendarAPITestSuite struct {
	suite.Suite
	db       *database.Database
	httpMock *internal.EndpointsMock
}

func TestCalendarAPITestSuite(t *testing.T) {

	suite.Run(t, new(CalendarAPITestSuite))
}

func (s *CalendarAPITestSuite) SetupTest() {
	s.httpMock = &internal.EndpointsMock{}
	managers.Microservices = s.httpMock
	_ = database.RemoveDB(databaseTest)
	s.db = database.InitDB(databaseTest)

	s.db.Conn.Exec("INSERT INTO calendar (meal_id,user_id,date,name) VALUES (?,?,?,?)", "01FN3EEB2NVFJAHAPM00000001", "01FN3EEB2NVFJAHAPU00000002", time.Now().Format("2006/01/02"), "pizza")
	s.db.Conn.Exec("INSERT INTO calendar (meal_id,user_id,date,name) VALUES (?,?,?,?)", "01FN3EEB2NVFJAHAPM00000002", "01FN3EEB2NVFJAHAPU00000002", time.Now().AddDate(0, 0, 1).Format("2006/01/02"), "pizza")

}

func (s *CalendarAPITestSuite) TearDownTest() {
	s.db = nil
	_ = database.RemoveDB(databaseTest)
}

func (s *CalendarAPITestSuite) TestPostCalendarHandler() {
	tests := []struct {
		name               string
		userId             string
		expectedULID       ulid.ULID
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Create new calendar (ok)",
			userId:             "01FN3EEB2NVFJAHAPU00000001",
			expectedStatusCode: http.StatusCreated,
			wantErr:            false,
		},
		{
			name: "User id not present (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, internal.RouteCalendar, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}

	for _, t := range tests {
		s.Run(t.name, func() {

			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}

			s.httpMock.On("GetAllMeals", t.userId).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userId, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}
			c := getEchoContext(t.userId)
			err := api.PostCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *CalendarAPITestSuite) TestGetCalendarHandler() {
	tests := []struct {
		name               string
		userID             string
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Get calendar (ok)",
			userID:             "01FN3EEB2NVFJAHAPU00000002",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Get calendar, userId not indicated (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "Get calendar, calendar not found (404)",
			userID: "01FN3EEB2NVFJAHAPU00000099",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusNotFound,
					Message: internal.ErrCalendarNotFound.Error(),
				},
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, internal.RouteCalendar, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}

			s.httpMock.On("GetAllMeals", t.userID).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userID, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}

			c := getEchoContext(t.userID)
			err := api.GetCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *CalendarAPITestSuite) TestPutCalendarHandler() {
	tests := []struct {
		name               string
		userID             string
		reqBody            interface{}
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "Update calendar (ok)",
			userID: "01FN3EEB2NVFJAHAPU00000002",
			reqBody: models.Calendar{
				MealId: "01FN3EEB2NVFJAHAPM00000010",
				Date:   time.Now().Format("2006/01/02"),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Update meal, userId not indicated (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string, request interface{}) echo.Context {
		var body []byte
		body, err := jsoniter.Marshal(request)
		s.NoError(err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, internal.RouteCalendar, bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}

			s.httpMock.On("GetAllMeals", t.userID).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userID, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}

			c := getEchoContext(t.userID, t.reqBody)
			err := api.PutCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}
			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *CalendarAPITestSuite) TestDeleteCalendarHandler() {
	tests := []struct {
		name               string
		userID             string
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Delete calendar (ok)",
			userID:             "01FN3EEB2NVFJAHAPU00000002",
			expectedStatusCode: http.StatusNoContent,
			wantErr:            false,
		},
		{
			name: "Delete calendar, userId not indicated (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:   "Calendar does not exist (404)",
			userID: "01FN3EEB2NVFJAHAPU00000002",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusNotFound,
					Message: internal.ErrCalendarNotFound.Error(),
				},
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, internal.RouteCalendar, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}
			s.httpMock.On("GetAllMeals", t.userID).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userID, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}
			c := getEchoContext(t.userID)
			err := api.DeleteCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}
			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *CalendarAPITestSuite) TestRedoCalendarHandler() {
	tests := []struct {
		name               string
		userID             string
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Redo calendar (ok)",
			userID:             "01FN3EEB2NVFJAHAPU00000002",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Redo calendar, userId not indicated (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, internal.RouteCalendarRedo, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}

			s.httpMock.On("GetAllMeals", t.userID).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userID, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}

			c := getEchoContext(t.userID)
			err := api.RedoCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}
			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *CalendarAPITestSuite) TestRedoWeekCalendarHandler() {
	tests := []struct {
		name               string
		userID             string
		reqBody            interface{}
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "Update calendar week (ok)",
			userID: "01FN3EEB2NVFJAHAPU00000002",
			reqBody: models.UpdateWeekCalendar{
				From: time.Now().Format("2006/01/02"),
				To:   time.Now().AddDate(0, 0, 1).Format("2006/01/02"),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "Update calendar week, userId not indicated (400)",
			expectedResp: &internal.ErrorResponse{
				Err: internal.ErrorBody{
					Status:  http.StatusBadRequest,
					Message: internal.ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string, request interface{}) echo.Context {
		var body []byte
		body, err := jsoniter.Marshal(request)
		s.NoError(err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, internal.RouteCalendarRedoWeek, bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(internal.ParamUserID)
		c.SetParamValues(userId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			calendarManager := managers.NewCalendarManager(*s.db)
			api := CalendarAPI{DB: *s.db, Manager: calendarManager}

			s.httpMock.On("GetAllMeals", t.userID).Return(mealsDb, nil).Once()
			for i, meal := range mealsDb {
				s.httpMock.On("GetMeal", t.userID, meal.Id).Return(models.MealToFront{Name: fmt.Sprintf("meal%d", i)}, nil)
			}

			c := getEchoContext(t.userID, t.reqBody)
			err := api.RedoWeekCalendarHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(internal.ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}
			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}
