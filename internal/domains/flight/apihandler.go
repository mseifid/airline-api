package flight

import (
	"errors"
	"net/http"
	"time"

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

// CreateFlight godoc
//
//	@Summary		Create a flight
//	@Description	Creates a new flight.
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			request	body		CreateFlightRequest	true	"Flight information"
//	@Success		201		{object}	api.APIResponse[FlightResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/internal/flights [post]
func (h *Handler) Create(c *echo.Context) error {
	var req CreateFlightRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	flight, err := h.service.Create(
		c.Request().Context(),
		req.AirplaneID,
		req.DepartureAirportID,
		req.ArrivalAirportID,
		req.DepartureAt,
		req.ArrivalAt,
		req.Price,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusCreated,
		api.Success(toResponse(&flight)),
	)
}

// UpdateFlight godoc
//
//	@Summary		Update a flight
//	@Description	Updates an existing flight.
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Security		AirlineAPIKey
//	@Param			id		path		int					true	"Flight ID"
//	@Param			request	body		UpdateFlightRequest	true	"Flight information"
//	@Success		200		{object}	api.APIResponse[FlightResponse]
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/internal/flights/{id} [patch]
func (h *Handler) Update(c *echo.Context) error {
	id, err := echo.PathParam[int64](c, "id")
	if err != nil || id <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid flight id",
		)
	}

	var req UpdateFlightRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	flight, err := h.service.Update(
		c.Request().Context(),
		id,
		req.AirplaneID,
		req.DepartureAirportID,
		req.ArrivalAirportID,
		req.DepartureAt,
		req.ArrivalAt,
		req.Price,
		req.Status,
	)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"flight not found",
			)
		}

		return echo.NewHTTPError(
			http.StatusBadRequest,
			err.Error(),
		)
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(&flight)),
	)
}

// Search flights
//
// @Summary Search flights
// @Tags Flights
// @Security AgencyAPIKey
// @Produce json
// @Param source_city query int true "Source city ID"
// @Param destination_city query int true "Destination city ID"
// @Param departure_date query string true "Start date (YYYY-MM-DD)"
// @Param arrival_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} api.APIResponse[[]FlightResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flights [get]
func (h *Handler) Search(c *echo.Context) error {
	sourceCityID, err := echo.QueryParam[int64](c, "source_city")
	if err != nil || sourceCityID <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"source_city is required",
		)
	}

	destinationCityID, err := echo.QueryParam[int64](c, "destination_city")
	if err != nil || destinationCityID <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"destination_city is required",
		)
	}

	departureDate := c.QueryParam("departure_date")
	if departureDate == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"departure_date is required",
		)
	}

	from, err := time.Parse("2006-01-02", departureDate)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid departure_date",
		)
	}

	// By default, search only the departure date.
	to := from.Add(24 * time.Hour)

	arrivalDate := c.QueryParam("arrival_date")

	if arrivalDate != "" {
		to, err = time.Parse("2006-01-02", arrivalDate)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid arrival_date",
			)
		}

		// Make the arrival date inclusive.
		to = to.Add(24 * time.Hour)

		if !to.After(from) {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"arrival_date must be on or after departure_date",
			)
		}
	}

	flights, err := h.service.Search(
		c.Request().Context(),
		sourceCityID,
		destinationCityID,
		from,
		to,
	)
	if err != nil {
		return err
	}

	response := make([]FlightResponse, 0, len(flights))

	for i := range flights {
		response = append(response, toResponse(&flights[i]))
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}

// List active flights
//
// @Summary List active flights
// @Tags Flights
// @Security AirlineAPIKey
// @Produce json
// @Success 200 {object} api.APIResponse[[]FlightResponse]
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /internal/flights [get]
func (h *Handler) ListActive(c *echo.Context) error {
	flights, err := h.service.ListActive(c.Request().Context())
	if err != nil {
		return err
	}

	response := make([]FlightResponse, 0, len(flights))

	for i := range flights {
		response = append(response, toResponse(&flights[i]))
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}

func toResponse(flight *Flight) FlightResponse {
	return FlightResponse{
		ID:                 flight.ID,
		AirplaneID:         flight.AirplaneID,
		DepartureAirportID: flight.DepartureAirportID,
		ArrivalAirportID:   flight.ArrivalAirportID,
		DepartureAt:        flight.DepartureAt,
		ArrivalAt:          flight.ArrivalAt,
		Price:              flight.Price,
		Capacity:           flight.Capacity,
		AvailableSeats:     flight.AvailableSeats,
		Status:             flight.Status,
	}
}