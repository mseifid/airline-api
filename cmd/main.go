package main

import (
	"errors"
	"firefly-airline/internal/api"
	"firefly-airline/internal/config"
	"firefly-airline/internal/domains/agency"
	"firefly-airline/internal/domains/airplane"
	"firefly-airline/internal/infrastructure/database"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "firefly-airline/api"

	"github.com/labstack/echo/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/swaggo/swag/v2"
)

// @title           Firefly Airline API
// @version         1.0
// @description     API documentation for Firefly Airline
//
//	@host			localhost:3000
//	@BasePath		/
//
// @securityDefinitions.apikey AirlineAPIKey
// @in header
// @name X-API-Key
//
// @securityDefinitions.apikey AgencyAPIKey
// @in header
// @name X-API-Key
func main() {
	// config loading
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config loaded successfully")

	// db
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	fmt.Println("database connection successful")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(sqlDB, "./migrations"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("migrations applied successfully")

	// Wiring
	agencyRepository := database.NewAgencyRepository(db)
	agencyService := agency.NewService(agencyRepository)
	agencyHandler := agency.NewHandler(agencyService)

	airplaneRepository := database.NewAirplaneRepository(db)
	airplaneService := airplane.NewService(airplaneRepository)
	airplaneHandler := airplane.NewHandler(airplaneService)

	// API server
	e := echo.New()

	// Global error handler
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		code := http.StatusInternalServerError
		message := "internal server error"

		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			code = httpErr.Code
			message = httpErr.Message
		}

		_ = c.JSON(code, api.Error(message))
	}

	internal := e.Group(
		"/internal",
		AirlineAPIKeyMiddleware(cfg.AirlineAPIKey),
	)

	internal.POST("/agencies", agencyHandler.Create)
	internal.POST("/airplanes", airplaneHandler.Create)
	internal.GET("/airplanes", airplaneHandler.List)
	internal.GET("/airplanes/:id", airplaneHandler.GetByID)
	internal.PATCH("/airplanes/:id", airplaneHandler.Update)
	internal.DELETE("/airplanes/:id", airplaneHandler.Delete)

	// Register the Swagger JSON endpoint
	e.GET("/swagger/doc.json", func(c *echo.Context) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read swagger doc"})
		}
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		return c.Blob(http.StatusOK, "application/json", []byte(doc))
	})

	// Register the Swagger UI route
	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	)))
	e.Start(":" + strconv.Itoa(cfg.Server.Port))
}

// TODO: move it to better place
func AirlineAPIKeyMiddleware(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")

			if apiKey == "" || apiKey != expectedKey {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid api key",
				)
			}

			return next(c)
		}
	}
}
