package country

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"firefly-airline/internal/api"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCountry godoc
//
//	@Summary		Create a country
//	@Description	Creates a new country.
//	@Tags			Countries
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateCountryRequest	true	"Country information"
//	@Success		201		{object}	api.APIResponse[CountryResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/countries [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateCountryRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	country, err := h.service.Create(
		c.Request().Context(),
		req.Name,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusCreated,
		api.Success(toResponse(&country)),
	)
}

// GetCountry godoc
//
//	@Summary		Get a country
//	@Description	Gets a country by ID.
//	@Tags			Countries
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id	path		int	true	"Country ID"
//	@Success		200	{object}	api.APIResponse[CountryResponse]
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/internal/countries/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid country id",
		)
	}

	country, err := h.service.GetByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"country not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&country)),
	)
}

// ListCountries godoc
//
//	@Summary		List countries
//	@Description	Returns all countries.
//	@Tags			Countries
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Success		200	{object}	api.APIResponse[[]CountryResponse]
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/internal/countries [get]
func (h *Handler) List(c *echo.Context) error {
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

// UpdateCountry godoc
//
//	@Summary		Update a country
//	@Description	Updates a country.
//	@Tags			Countries
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int					true	"Country ID"
//	@Param			request	body		UpdateCountryRequest	true	"Country information"
//	@Success		200		{object}	api.APIResponse[CountryResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/countries/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid country id",
		)
	}

	var req UpdateCountryRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Name = strings.TrimSpace(req.Name)

	country, err := h.service.Update(
		c.Request().Context(),
		id,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"country not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&country)),
	)
}

func toResponse(country *Country) *CountryResponse {
	return &CountryResponse{
		ID:   country.ID,
		Name: country.Name,
	}
}