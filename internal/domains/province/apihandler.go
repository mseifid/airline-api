package province

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

// CreateProvince godoc
//
//	@Summary		Create a province
//	@Description	Creates a new province.
//	@Tags			Provinces
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateProvinceRequest	true	"Province information"
//	@Success		201		{object}	api.APIResponse[ProvinceResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/provinces [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateProvinceRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	province, err := h.service.Create(
		c.Request().Context(),
		req.CountryID,
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
		api.Success(toResponse(&province)),
	)
}

// GetProvince godoc
//
//	@Summary		Get a province
//	@Description	Gets a province by ID.
//	@Tags			Provinces
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id	path		int	true	"Province ID"
//	@Success		200	{object}	api.APIResponse[ProvinceResponse]
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/internal/provinces/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid province id",
		)
	}

	province, err := h.service.GetByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"province not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&province)),
	)
}

// ListProvinces godoc
//
//	@Summary		List provinces
//	@Description	Returns all provinces, optionally filtered by country.
//	@Tags			Provinces
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			country_id	query		int	false	"Filter by country ID"
//	@Success		200			{object}	api.APIResponse[[]ProvinceResponse]
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/internal/provinces [get]
func (h *Handler) List(c *echo.Context) error {
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

// UpdateProvince godoc
//
//	@Summary		Update a province
//	@Description	Updates a province.
//	@Tags			Provinces
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int					true	"Province ID"
//	@Param			request	body		UpdateProvinceRequest	true	"Province information"
//	@Success		200		{object}	api.APIResponse[ProvinceResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/provinces/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid province id",
		)
	}

	var req UpdateProvinceRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	province, err := h.service.Update(
		c.Request().Context(),
		id,
		req.CountryID,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"province not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&province)),
	)
}

func toResponse(province *Province) *ProvinceResponse {
	return &ProvinceResponse{
		ID:        province.ID,
		CountryID: province.CountryID,
		Name:      province.Name,
	}
}