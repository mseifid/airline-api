package country

import (
	"net/http"

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

// ListCountries godoc
//
//	@Summary		List countries
//	@Description	Returns all countries.
//	@Tags			Countries
//	@Produce		json
//	@Security		AgencyAPIKey
//	@Success		200	{object}	api.APIResponse[[]CountryResponse]
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/countries [get]
func (h *PublicHandler) List(c *echo.Context) error {
	countries, err := h.service.List(
		c.Request().Context(),
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	response := make([]CountryResponse, 0, len(countries))

	for i := range countries {
		response = append(
			response,
			*toResponse(&countries[i]),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}