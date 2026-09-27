package flight

import "time"

type Flight struct {
	ID                  int64
	AirplaneID          int64
	DepartureAirportID  int64
	ArrivalAirportID    int64
	DepartureAt         time.Time
	ArrivalAt           time.Time
	Price               int64
	Capacity            int
	AvailableSeats      int
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}