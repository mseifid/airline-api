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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config loaded successfully")

	_, err = database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	fmt.Println("database connection successful")

	e := echo.New()

	e.GET("/sample", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Sample endpoint",
		})
	})

	e.Start(":" + strconv.Itoa(cfg.Server.Port))
}
