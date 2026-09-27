package airplane

import (
	"firefly-airline/internal/api"
	"net/http"
	"strconv"
	"strings"

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

// CreateAirplane godoc
//
//	@Summary		Create an airplane
//	@Description	Creates a new airplane.
//	@Tags			Airplanes
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateAirplaneRequest	true	"Airplane information"
//	@Success		201		{object}	AirplaneResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/internal/airplanes [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateAirplaneRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Type = strings.TrimSpace(req.Type)

	airplane, err := h.service.Create(
		c.Request().Context(),
		req.Type,
		req.Capacity,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusCreated,
		api.Success(toResponse(airplane)),
	)
}

// GetAirplane godoc
//
//	@Summary		Get an airplane
//	@Description	Gets an airplane by ID.
//	@Tags			Airplanes
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id	path		int	true	"Airplane ID"
//	@Success		200	{object}	AirplaneResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/internal/airplanes/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid airplane id",
		)
	}

	airplane, err := h.service.GetByID(
		c.Request().Context(),
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(
				http.StatusNotFound,
				err.Error(),
			)
		}

		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(airplane)),
	)
}

// ListAirplanes godoc
//
//	@Summary		List airplanes
//	@Description	Returns all airplanes.
//	@Tags			Airplanes
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Success		200	{array}	AirplaneResponse
//	@Failure		401	{object}	map[string]string
//	@Router			/internal/airplanes [get]
func (h *Handler) List(c *echo.Context) error {
	airplanes, err := h.service.List(
		c.Request().Context(),
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	response := make([]AirplaneResponse, 0, len(airplanes))

	for i := range airplanes {
		response = append(response, *toResponse(&airplanes[i]))
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}

// UpdateAirplane godoc
//
//	@Summary		Update an airplane
//	@Description	Updates an airplane.
//	@Tags			Airplanes
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int						true	"Airplane ID"
//	@Param			request	body		UpdateAirplaneRequest	true	"Airplane information"
//	@Success		200		{object}	AirplaneResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/internal/airplanes/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid airplane id",
		)
	}

	var req UpdateAirplaneRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	req.Type = strings.TrimSpace(req.Type)

	airplane, err := h.service.Update(
		c.Request().Context(),
		id,
		req.Type,
		req.Capacity,
	)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(
				http.StatusNotFound,
				err.Error(),
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	// The update repository doesn't currently return
	// the updated timestamps, so fetch the updated entity.
	airplane, err = h.service.GetByID(
		c.Request().Context(),
		airplane.ID,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(airplane)),
	)
}

// DeleteAirplane godoc
//
//	@Summary		Delete an airplane
//	@Description	Deletes an airplane. Deletion fails if the airplane is referenced by a flight.
//	@Tags			Airplanes
//	@Security		AirlineAPIKey
//	@Param			id	path	int	true	"Airplane ID"
//	@Success		204
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		409	{object}	map[string]string
//	@Router			/internal/airplanes/{id} [delete]
func (h *Handler) Delete(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid airplane id",
		)
	}

	if err := h.service.Delete(
		c.Request().Context(),
		id,
	); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(
				http.StatusNotFound,
				err.Error(),
			)
		}

		return echo.NewHTTPError(
			http.StatusConflict,
			err.Error(),
		)
	}

	return c.NoContent(http.StatusNoContent)
}

func toResponse(a *Airplane) *AirplaneResponse {
	return &AirplaneResponse{
		ID:       a.ID,
		Type:     a.Type,
		Capacity: a.Capacity,
	}
}
