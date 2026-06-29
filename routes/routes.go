package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/handler"
	"github.com/shihabul3000/Go-Project/middleware"
	"github.com/shihabul3000/Go-Project/models"
)

func Register(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler,
	jwtSecret string,
) {
	api := e.Group("/api/v1")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	api.GET("/zones", zoneHandler.List)
	api.GET("/zones/:id", zoneHandler.Get)
	api.POST("/zones", zoneHandler.Create, middleware.JWTAuth(jwtSecret), middleware.RequireRoles(models.RoleAdmin))
	api.PATCH("/zones/:id", zoneHandler.Update, middleware.JWTAuth(jwtSecret), middleware.RequireRoles(models.RoleAdmin))
	api.DELETE("/zones/:id", zoneHandler.Delete, middleware.JWTAuth(jwtSecret), middleware.RequireRoles(models.RoleAdmin))

	reservations := api.Group("/reservations", middleware.JWTAuth(jwtSecret))
	reservations.POST("", reservationHandler.Create)
	reservations.GET("/my-reservations", reservationHandler.MyReservations)
	reservations.DELETE("/:id", reservationHandler.Cancel)
	reservations.GET("", reservationHandler.ListAll, middleware.RequireRoles(models.RoleAdmin))
}
