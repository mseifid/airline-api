package city

import (
	"net/http"
	"strconv"

	"firefly-airline/internal/api"

	"github.com/labstack/echo/v5"
)

type PublicHandler struct {
	service *Service
}

func NewPublicHandler(service *Service) *PublicHandler {
	return &PublicHandler{
		service: service,
	}
}

// ListCities godoc
//
//	@Summary		List cities
//	@Description	Returns cities, optionally filtered by province.
//	@Tags			Cities
//	@Produce		json
//	@Security		AgencyAPIKey
//	@Param			province_id	query		int	false	"Filter by province ID"
//	@Success		200			{object}	api.APIResponse[[]CityResponse]
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/cities [get]
func (h *PublicHandler) List(c *echo.Context) error {
	var provinceID *int64

	value := c.QueryParam("province_id")

	if value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid province id",
			)
		}

		provinceID = &id
	}

	cities, err := h.service.List(
		c.Request().Context(),
		provinceID,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	response := make([]CityResponse, 0, len(cities))

	for i := range cities {
		response = append(
			response,
			*toResponse(&cities[i]),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}