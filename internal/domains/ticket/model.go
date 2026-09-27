package ticket

import "time"

const (
	StatusActive    = "active"
	StatusCancelled = "cancelled"
)

type Ticket struct {
	ID           int64
	AgencyID     int64
	FlightID     int64
	SeatCount    int
	UnitPrice    int64
	TotalPrice   int64
	Status       string
	CreatedAt    time.Time
	CancelledAt  *time.Time

	DepartureAt  time.Time
	ArrivalAt    time.Time
	AirplaneType string

	Passengers []Passenger
}

type Passenger struct {
	ID          int64
	TicketID    int64
	Name        string
	Mobile      string
	NationalCode string
}