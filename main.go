package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/shihabul3000/Go-Project/config"
	"github.com/shihabul3000/Go-Project/database"
	"github.com/shihabul3000/Go-Project/handler"
	"github.com/shihabul3000/Go-Project/repository"
	"github.com/shihabul3000/Go-Project/routes"
	"github.com/shihabul3000/Go-Project/service"
	"github.com/shihabul3000/Go-Project/validation"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiresIn, cfg.BcryptCost)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo)

	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	e := echo.New()
	e.Validator = validation.New()
	e.HTTPErrorHandler = handler.HTTPErrorHandler

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	routes.Register(e, authHandler, zoneHandler, reservationHandler, cfg.JWTSecret)

	log.Printf("SpotSync API is running on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
