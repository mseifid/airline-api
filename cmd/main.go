package main

import (
	"errors"
	"firefly-airline/internal/api"
	"firefly-airline/internal/config"
	"firefly-airline/internal/domains/agency"
	"firefly-airline/internal/domains/airplane"
	"firefly-airline/internal/domains/airport"
	"firefly-airline/internal/domains/city"
	"firefly-airline/internal/domains/country"
	"firefly-airline/internal/domains/province"
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

	countryRepository := database.NewCountryRepository(db)
	countryService := country.NewService(countryRepository)
	countryHandler := country.NewHandler(countryService)
	countryPublicHandler := country.NewPublicHandler(countryService)

	provinceRepository := database.NewProvinceRepository(db)
	provinceService := province.NewService(provinceRepository, countryRepository)
	provinceHandler := province.NewHandler(provinceService)
	provincePublicHandler := province.NewPublicHandler(provinceService)

	cityRepository := database.NewCityRepository(db)
	cityService := city.NewService(cityRepository, provinceRepository)
	cityHandler := city.NewHandler(cityService)
	cityPublicHandler := city.NewPublicHandler(cityService)

	airportRepository := database.NewAirportRepository(db)
	airportService := airport.NewService(airportRepository, cityRepository)
	airportHandler := airport.NewHandler(airportService)
	airportPublicHandler := airport.NewPublicHandler(airportService)

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
		api.AirlineAPIKeyMiddleware(cfg.AirlineAPIKey),
	)

	internal.POST("/agencies", agencyHandler.Create)

	internal.POST("/airplanes", airplaneHandler.Create)
	internal.GET("/airplanes", airplaneHandler.List)
	internal.GET("/airplanes/:id", airplaneHandler.GetByID)
	internal.PATCH("/airplanes/:id", airplaneHandler.Update)
	internal.DELETE("/airplanes/:id", airplaneHandler.Delete)

	internal.POST("/countries", countryHandler.Create)
	internal.GET("/countries", countryHandler.List)
	internal.GET("/countries/:id", countryHandler.GetByID)
	internal.PATCH("/countries/:id", countryHandler.Update)

	internal.POST("/provinces", provinceHandler.Create)
	internal.GET("/provinces", provinceHandler.List)
	internal.GET("/provinces/:id", provinceHandler.GetByID)
	internal.PATCH("/provinces/:id", provinceHandler.Update)

	internal.POST("/cities", cityHandler.Create)
	internal.GET("/cities", cityHandler.List)
	internal.GET("/cities/:id", cityHandler.GetByID)
	internal.PATCH("/cities/:id", cityHandler.Update)

	internal.POST("/airports", airportHandler.Create)
	internal.GET("/airports", airportHandler.List)
	internal.GET("/airports/:id", airportHandler.GetByID)
	internal.PATCH("/airports/:id", airportHandler.Update)

	public := e.Group(
		"",
		api.AgencyAPIKeyMiddleware(agencyService),
	)

	public.GET("/countries", countryPublicHandler.List)
	public.GET("/provinces", provincePublicHandler.List)
	public.GET("/cities", cityPublicHandler.List)
	public.GET("/airports", airportPublicHandler.List)

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
