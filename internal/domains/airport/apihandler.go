package airport

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

// CreateAirport godoc
//
//	@Summary		Create an airport
//	@Description	Creates a new airport.
//	@Tags			Airports
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateAirportRequest	true	"Airport information"
//	@Success		201		{object}	api.APIResponse[AirportResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/internal/airports [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateAirportRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	airport, err := h.service.Create(
		c.Request().Context(),
		req.CityID,
		req.Code,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			return echo.NewHTTPError(
				http.StatusConflict,
				"airport code already exists",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusCreated,
		api.Success(toResponse(&airport)),
	)
}

// GetAirport godoc
//
//	@Summary		Get an airport
//	@Description	Gets an airport by ID.
//	@Tags			Airports
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id	path		int	true	"Airport ID"
//	@Success		200	{object}	api.APIResponse[AirportResponse]
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/internal/airports/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid airport id",
		)
	}

	airport, err := h.service.GetByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"airport not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&airport)),
	)
}

// ListAirports godoc
//
//	@Summary		List airports
//	@Description	Returns all airports, optionally filtered by city.
//	@Tags			Airports
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			city_id	query		int	false	"Filter by city ID"
//	@Success		200		{object}	api.APIResponse[[]AirportResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/internal/airports [get]
func (h *Handler) List(c *echo.Context) error {
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

	response := make(
		[]AirportResponse,
		0,
		len(airports),
	)

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

// UpdateAirport godoc
//
//	@Summary		Update an airport
//	@Description	Updates an airport.
//	@Tags			Airports
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int					true	"Airport ID"
//	@Param			request	body		UpdateAirportRequest	true	"Airport information"
//	@Success		200		{object}	api.APIResponse[AirportResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Router			/internal/airports/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid airport id",
		)
	}

	var req UpdateAirportRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)

	airport, err := h.service.Update(
		c.Request().Context(),
		id,
		req.CityID,
		req.Code,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"airport not found",
			)
		}

		if errors.Is(err, ErrCodeAlreadyExists) {
			return echo.NewHTTPError(
				http.StatusConflict,
				"airport code already exists",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&airport)),
	)
}

func toResponse(airport *Airport) *AirportResponse {
	return &AirportResponse{
		ID:     airport.ID,
		CityID: airport.CityID,
		Code:   airport.Code,
		Name:   airport.Name,
	}
}