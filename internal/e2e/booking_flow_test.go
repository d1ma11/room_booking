package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"test-backend-1-d1ma11/configs"
	"test-backend-1-d1ma11/internal/controller"
	"test-backend-1-d1ma11/internal/middleware"
	"test-backend-1-d1ma11/internal/repository"
	"test-backend-1-d1ma11/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BookingE2ESuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *pg.PostgresContainer
	dbDSN       string
	router      *gin.Engine
	cfg         *configs.Config

	adminToken string
	userToken  string
	roomID     string
	slotID     string
	bookingID  string
}

func (s *BookingE2ESuite) SetupSuite() {
	s.ctx = context.Background()
	gin.SetMode(gin.TestMode)

	pgContainer, err := pg.Run(s.ctx,
		"postgres:15-alpine",
		pg.WithDatabase("room_booking"),
		pg.WithUsername("postgres"),
		pg.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	s.Require().NoError(err, "Failed to start postgres container")
	s.pgContainer = pgContainer

	dsn, err := pgContainer.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err)
	s.dbDSN = dsn

	os.Setenv("PG_URL", dsn)

	m, err := migrate.New("file://../../migrations", dsn)
	s.Require().NoError(err, "Failed to init migrate")
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.FailNow("Migration failed", err.Error())
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err, "Failed to connect to test DB")

	s.cfg = &configs.Config{
		JWT: configs.JWTConfig{
			Secret: "test-secret-e2e",
			TTL:    time.Hour,
		},
	}

	repos := repository.NewRepositories(db)
	deps := service.ServicesDependencies{
		Repos:             repos,
		ConferenceService: service.NewConferenceService(),
		Config:            s.cfg,
	}
	services := service.NewServices(deps)

	s.router = gin.Default()
	s.router.Use(gin.Recovery())

	authCtrl := controller.NewUserController(services.UserService)
	s.router.POST("/dummyLogin", authCtrl.DummyLogin)

	authorized := s.router.Group("/")
	authorized.Use(middleware.AuthMiddleware(services.JwtService, s.cfg.JWT.Secret))

	roomCtrl := controller.NewRoomController(services.RoomService)
	authorized.GET("/rooms/list", roomCtrl.List)
	authorized.POST("/rooms/create", middleware.RequireRole("admin"), roomCtrl.Create)

	scheduleCtrl := controller.NewScheduleController(services.ScheduleService)
	authorized.POST("/rooms/:roomId/schedule/create", middleware.RequireRole("admin"), scheduleCtrl.Create)

	slotCtrl := controller.NewSlotController(services.SlotService)
	authorized.GET("/rooms/:roomId/slots/list", slotCtrl.List)

	bookingCtrl := controller.NewBookingController(services.BookingService)
	authorized.POST("/bookings/create", middleware.RequireRole("user"), bookingCtrl.BookSlot)
	authorized.POST("/bookings/:bookingId/cancel", middleware.RequireRole("user"), bookingCtrl.CancelBooking)
	authorized.GET("/bookings/my", middleware.RequireRole("user"), bookingCtrl.ListMy)
	authorized.GET("/bookings/list", middleware.RequireRole("admin"), bookingCtrl.ListAdmin)

	s.router.GET("/_info", func(c *gin.Context) { c.Status(http.StatusOK) })

	s.adminToken = s.generateToken(services.JwtService, "00000000-0000-0000-0000-000000000001", "admin")
	s.userToken = s.generateToken(services.JwtService, "00000000-0000-0000-0000-000000000002", "user")
}

func (s *BookingE2ESuite) TearDownSuite() {
	if s.pgContainer != nil {
		err := s.pgContainer.Terminate(s.ctx)
		s.Require().NoError(err, "Failed to terminate container")
	}
}

func (s *BookingE2ESuite) generateToken(jwtService service.JwtService, userID, role string) string {
	token, err := jwtService.GenerateToken(userID, role, s.cfg.JWT.Secret, s.cfg.JWT.TTL)
	s.Require().NoError(err)
	return token
}

func (s *BookingE2ESuite) performRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

func (s *BookingE2ESuite) TestE2E_FullBookingFlow() {
	// Создаем переговорку через админку
	roomReq := map[string]interface{}{
		"name":        "E2E Test Room",
		"description": "Room for E2E tests",
		"capacity":    10,
	}
	w := s.performRequest("POST", "/rooms/create", roomReq, s.adminToken)
	s.Equal(http.StatusCreated, w.Code, "Failed to create room")

	var roomResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	s.roomID = roomResp["room"].(map[string]interface{})["id"].(string)
	s.NotEmpty(s.roomID)

	// Создаем расписание через админку
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	dayOfWeek := int(tomorrow.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	scheduleReq := map[string]interface{}{
		"daysOfWeek": []int{dayOfWeek},
		"startTime":  "09:00",
		"endTime":    "11:00",
	}
	w = s.performRequest("POST", fmt.Sprintf("/rooms/%s/schedule/create", s.roomID), scheduleReq, s.adminToken)
	s.Equal(http.StatusCreated, w.Code, "Failed to create schedule: %s", w.Body.String())

	// Получаем список всех доступных слотов через пользователя
	dateStr := tomorrow.Format("2006-01-02")
	w = s.performRequest("GET", fmt.Sprintf("/rooms/%s/slots/list?date=%s", s.roomID, dateStr), nil, s.userToken)
	s.Equal(http.StatusOK, w.Code)

	var slotsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &slotsResp)
	slots := slotsResp["slots"].([]interface{})
	s.GreaterOrEqual(len(slots), 1, "No slots generated")

	// Берем первый доступный слот
	firstSlot := slots[0].(map[string]interface{})
	s.slotID = firstSlot["id"].(string)

	// Делаем бронь через пользователя
	bookingReq := map[string]interface{}{
		"slotId":               s.slotID,
		"createConferenceLink": true,
	}
	w = s.performRequest("POST", "/bookings/create", bookingReq, s.userToken)
	s.Equal(http.StatusCreated, w.Code, "Failed to create booking: %s", w.Body.String())

	var bookingResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &bookingResp)
	bookingData := bookingResp["booking"].(map[string]interface{})
	s.bookingID = bookingData["id"].(string)
	s.Equal("active", bookingData["status"])

	// Проверяем, что создалась ссылка на конференцию
	confLink, exists := bookingData["conferenceLink"]
	s.True(exists, "Conference link should exist")
	s.Contains(confLink, "https://conference.com/meeting/")
}

func (s *BookingE2ESuite) TestE2E_CancelBooking() {
	// Для начала создадим переговорку
	roomReq := map[string]interface{}{"name": "Cancel Test Room"}
	w := s.performRequest("POST", "/rooms/create", roomReq, s.adminToken)
	s.Equal(http.StatusCreated, w.Code)
	var roomResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &roomResp)
	roomID := roomResp["room"].(map[string]interface{})["id"].(string)

	// И добавим расписание к этой переговорке
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	dayOfWeek := int(tomorrow.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	scheduleReq := map[string]interface{}{
		"daysOfWeek": []int{dayOfWeek},
		"startTime":  "12:00",
		"endTime":    "13:00",
	}
	w = s.performRequest("POST", fmt.Sprintf("/rooms/%s/schedule/create", roomID), scheduleReq, s.adminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Получаем слоты для этой переговорки
	w = s.performRequest("GET", fmt.Sprintf("/rooms/%s/slots/list?date=%s", roomID, tomorrow.Format("2006-01-02")), nil, s.userToken)
	s.Equal(http.StatusOK, w.Code)
	var slotsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &slotsResp)
	slots := slotsResp["slots"].([]interface{})
	s.Greater(len(slots), 0)
	slotID := slots[0].(map[string]interface{})["id"].(string)

	// Бронируем слот
	bookingReq := map[string]interface{}{"slotId": slotID, "createConferenceLink": false}
	w = s.performRequest("POST", "/bookings/create", bookingReq, s.userToken)
	s.Equal(http.StatusCreated, w.Code)
	var bookingResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &bookingResp)
	bookingID := bookingResp["booking"].(map[string]interface{})["id"].(string)

	// И отменяем этот слот
	w = s.performRequest("POST", fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, s.userToken)
	s.Equal(http.StatusOK, w.Code, "Failed to cancel booking")

	var cancelResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cancelResp)
	status := cancelResp["booking"].(map[string]interface{})["status"]
	s.Equal("cancelled", status)

	// Проверим, что при повторной отмене возвращается статус ОК 200
	w = s.performRequest("POST", fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, s.userToken)
	s.Equal(http.StatusOK, w.Code, "Re-canceling should be idempotent")
}

func TestBookingE2ESuite(t *testing.T) {
	suite.Run(t, new(BookingE2ESuite))
}
