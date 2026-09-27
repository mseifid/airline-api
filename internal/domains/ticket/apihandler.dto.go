package ticket

import "time"

type PassengerRequest struct {
	Name        string `json:"name" example:"John Doe"`
	Mobile      string `json:"mobile" example:"+989121234567"`
	NationalCode string `json:"nationalCode" example:"0012345678"`
}

type PurchaseTicketRequest struct {
	FlightID  int64             `json:"flightId" example:"10"`
	Passengers []PassengerRequest `json:"passengers"`
}

type PassengerResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Mobile       string `json:"mobile"`
	NationalCode string `json:"nationalCode"`
}

type TicketResponse struct {
	ID           int64              `json:"id"`
	AgencyID     int64              `json:"agencyId"`
	FlightID     int64              `json:"flightId"`
	SeatCount    int                `json:"seatCount"`
	UnitPrice    int64              `json:"unitPrice"`
	TotalPrice   int64              `json:"totalPrice"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"createdAt"`
	CancelledAt  *time.Time         `json:"cancelledAt,omitempty"`
	DepartureAt  time.Time          `json:"departureAt"`
	ArrivalAt    time.Time          `json:"arrivalAt"`
	AirplaneType string             `json:"airplaneType"`
	Passengers   []PassengerResponse `json:"passengers"`
}