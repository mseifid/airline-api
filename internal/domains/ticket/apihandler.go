package ticket

import (
	"net/http"

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

// Purchase ticket
//
// @Summary Purchase ticket
// @Tags Tickets
// @Security AgencyAPIKey
// @Accept json
// @Produce json
// @Param request body PurchaseTicketRequest true "Ticket purchase data"
// @Success 201 {object} api.APIResponse[TicketResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets [post]
func (h *Handler) Purchase(c *echo.Context) error {
	var req PurchaseTicketRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	agencyID, ok := c.Get("agency_id").(int64)
	if !ok || agencyID <= 0 {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid agency context",
		)
	}

	passengers := make([]Passenger, 0, len(req.Passengers))

	for _, p := range req.Passengers {
		passengers = append(passengers, Passenger{
			Name:         p.Name,
			Mobile:       p.Mobile,
			NationalCode: p.NationalCode,
		})
	}

	ticket, err := h.service.Purchase(
		c.Request().Context(),
		agencyID,
		req.FlightID,
		passengers,
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusCreated,
		api.Success(toResponse(ticket)),
	)
}

// List tickets purchased by the agency for a flight
//
// @Summary List agency tickets for a flight
// @Tags Tickets
// @Security AgencyAPIKey
// @Produce json
// @Param flight_id path int true "Flight ID"
// @Success 200 {object} api.APIResponse[[]TicketResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flights/{flight_id}/tickets [get]
func (h *Handler) ListByFlight(c *echo.Context) error {
	flightID, err := echo.PathParam[int64](c, "flight_id")
	if err != nil || flightID <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid flight id",
		)
	}

	agencyID, ok := c.Get("agency_id").(int64)
	if !ok || agencyID <= 0 {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid agency context",
		)
	}

	tickets, err := h.service.ListByFlight(
		c.Request().Context(),
		agencyID,
		flightID,
	)
	if err != nil {
		return err
	}

	response := make([]TicketResponse, 0, len(tickets))

	for i := range tickets {
		response = append(response, toResponse(&tickets[i]))
	}

	return c.JSON(
		http.StatusOK,
		api.Success(response),
	)
}

// Cancel ticket
//
// @Summary Cancel ticket
// @Tags Tickets
// @Security AgencyAPIKey
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} api.APIResponse[TicketResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets/{id}/cancel [post]
func (h *Handler) Cancel(c *echo.Context) error {
	ticketID, err := echo.PathParam[int64](c, "id")
	if err != nil || ticketID <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid ticket id",
		)
	}

	agencyID, ok := c.Get("agency_id").(int64)
	if !ok || agencyID <= 0 {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid agency context",
		)
	}

	ticket, err := h.service.Cancel(
		c.Request().Context(),
		agencyID,
		ticketID,
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(ticket)),
	)
}

// Get ticket
//
// @Summary Get ticket
// @Tags Tickets
// @Security AgencyAPIKey
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} api.APIResponse[TicketResponse]
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tickets/{id} [get]
func (h *Handler) GetByID(c *echo.Context) error {
	ticketID, err := echo.PathParam[int64](c, "id")
	if err != nil || ticketID <= 0 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid ticket id",
		)
	}

	agencyID, ok := c.Get("agency_id").(int64)
	if !ok || agencyID <= 0 {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid agency context",
		)
	}

	result, err := h.service.GetByID(
		c.Request().Context(),
		ticketID,
		agencyID,
	)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		api.Success(toResponse(result)),
	)
}

func toResponse(ticket *Ticket) TicketResponse {
	passengers := make([]PassengerResponse, 0, len(ticket.Passengers))

	for i := range ticket.Passengers {
		passengers = append(passengers, PassengerResponse{
			ID:           ticket.Passengers[i].ID,
			Name:         ticket.Passengers[i].Name,
			Mobile:       ticket.Passengers[i].Mobile,
			NationalCode: ticket.Passengers[i].NationalCode,
		})
	}

	return TicketResponse{
		ID:           ticket.ID,
		AgencyID:     ticket.AgencyID,
		FlightID:     ticket.FlightID,
		SeatCount:    ticket.SeatCount,
		UnitPrice:    ticket.UnitPrice,
		TotalPrice:   ticket.TotalPrice,
		Status:       ticket.Status,
		CreatedAt:    ticket.CreatedAt,
		CancelledAt:  ticket.CancelledAt,
		DepartureAt:  ticket.DepartureAt,
		ArrivalAt:    ticket.ArrivalAt,
		AirplaneType: ticket.AirplaneType,
		Passengers:   passengers,
	}
}