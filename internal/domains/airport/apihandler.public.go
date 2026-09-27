package airport

import (
	"net/http"
	"strconv"

	"firefly-airline/internal/api"

	"github.com/labstack/echo/v5"
)

// PublicHandler handles agency-facing read-only airport endpoints.
type PublicHandler struct {
	service *Service
}

func NewPublicHandler(service *Service) *PublicHandler {
	return &PublicHandler{
		service: service,
	}
}

// ListAirports godoc
//
//	@Summary		List airports
//	@Description	Returns airports, optionally filtered by city.
//	@Tags			Airports
//	@Produce		json
//	@Security		AgencyAPIKey
//	@Param			city_id	query		int	false	"Filter by city ID"
//	@Success		200		{object}	api.APIResponse[[]AirportResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/airports [get]
func (h *PublicHandler) List(c *echo.Context) error {
	var cityID *int64

	value := c.QueryParam("city_id")

	if value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid city id",
			)
		}

		cityID = &id
	}

	airports, err := h.service.List(
		c.Request().Context(),
		cityID,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	response := make([]AirportResponse, 0, len(airports))

	for i := range airports {
		response = append(
			response,
			*toResponse(&airports[i]),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}