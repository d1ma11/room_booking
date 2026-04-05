package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-backend-1-d1ma11/configs"
	"test-backend-1-d1ma11/initializers"
	"test-backend-1-d1ma11/internal/controller"
	"test-backend-1-d1ma11/internal/middleware"
	"test-backend-1-d1ma11/internal/repository"
	"test-backend-1-d1ma11/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	if err = RunMigrations(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	db, err := initializers.ConnectToDb(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Repositories
	log.Println("Initializing repositories...")
	repositories := repository.NewRepositories(db)

	// Services dependencies
	log.Println("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:             repositories,
		ConferenceService: service.NewConferenceService(),
		Config:            cfg,
	}
	services := service.NewServices(deps)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/_info", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// NO AUTHORIZATION REQUIRED
	userCtrl := controller.NewUserController(services.UserService)
	router.POST("/dummyLogin", userCtrl.DummyLogin)
	router.POST("/register", userCtrl.Register)
	router.POST("/login", userCtrl.Login)

	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware(services.JwtService, cfg.JWT.Secret))

	// AUTHORIZED REQUIRED
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

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
