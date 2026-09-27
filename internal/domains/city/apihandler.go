package city

import (
	"errors"
	"net/http"
	"strconv"

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

// CreateCity godoc
//
//	@Summary		Create a city
//	@Description	Creates a new city.
//	@Tags			Cities
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateCityRequest	true	"City information"
//	@Success		201		{object}	api.APIResponse[CityResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/cities [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateCityRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	city, err := h.service.Create(
		c.Request().Context(),
		req.ProvinceID,
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
		api.Success(toResponse(&city)),
	)
}

// GetCity godoc
//
//	@Summary		Get a city
//	@Description	Gets a city by ID.
//	@Tags			Cities
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id	path		int	true	"City ID"
//	@Success		200	{object}	api.APIResponse[CityResponse]
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/internal/cities/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid city id",
		)
	}

	city, err := h.service.GetByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"city not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&city)),
	)
}

// ListCities godoc
//
//	@Summary		List cities
//	@Description	Returns all cities, optionally filtered by province.
//	@Tags			Cities
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			province_id	query		int	false	"Filter by province ID"
//	@Success		200			{object}	api.APIResponse[[]CityResponse]
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/internal/cities [get]
func (h *Handler) List(c *echo.Context) error {
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

// UpdateCity godoc
//
//	@Summary		Update a city
//	@Description	Updates a city.
//	@Tags			Cities
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int					true	"City ID"
//	@Param			request	body		UpdateCityRequest	true	"City information"
//	@Success		200		{object}	api.APIResponse[CityResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/cities/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid city id",
		)
	}

	var req UpdateCityRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	city, err := h.service.Update(
		c.Request().Context(),
		id,
		req.ProvinceID,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"city not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&city)),
	)
}

func toResponse(city *City) *CityResponse {
	return &CityResponse{
		ID:         city.ID,
		ProvinceID: city.ProvinceID,
		Name:       city.Name,
	}
}