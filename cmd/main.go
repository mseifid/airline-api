package main

import (
	"firefly-airline/internal/config"
	"firefly-airline/internal/infrastructure/database"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

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

	// API server
	e := echo.New()

	e.GET("/sample", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Sample endpoint",
		})
	})

	e.Start(":" + strconv.Itoa(cfg.Server.Port))
}
