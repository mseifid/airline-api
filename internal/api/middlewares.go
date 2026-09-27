package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AgencyAPIKeyValidator interface {
	ValidateAPIKey(ctx context.Context, rawKey string) (int64, error)
}

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

func AgencyAPIKeyMiddleware(validator AgencyAPIKeyValidator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")

			if apiKey == "" {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"missing api key",
				)
			}

			agencyID, err := validator.ValidateAPIKey(
				c.Request().Context(),
				apiKey,
			)
			if err != nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid api key",
				)
			}

			c.Set("agency_id", agencyID)

			return next(c)
		}
	}
}