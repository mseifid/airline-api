package province

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

// ListProvinces godoc
//
//	@Summary		List provinces
//	@Description	Returns provinces, optionally filtered by country.
//	@Tags			Provinces
//	@Produce		json
//	@Security		AgencyAPIKey
//	@Param			country_id	query		int	false	"Filter by country ID"
//	@Success		200			{object}	api.APIResponse[[]ProvinceResponse]
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/provinces [get]
func (h *PublicHandler) List(c *echo.Context) error {
	var countryID *int64

	value := c.QueryParam("country_id")

	if value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid country id",
			)
		}

		countryID = &id
	}

	provinces, err := h.service.List(
		c.Request().Context(),
		countryID,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	response := make([]ProvinceResponse, 0, len(provinces))

	for i := range provinces {
		response = append(
			response,
			*toResponse(&provinces[i]),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}