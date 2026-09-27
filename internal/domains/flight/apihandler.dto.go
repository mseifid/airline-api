package flight

import "time"

type CreateFlightRequest struct {
	AirplaneID         int64     `json:"airplaneId" example:"1"`
	DepartureAirportID int64     `json:"departureAirportId" example:"1"`
	ArrivalAirportID   int64     `json:"arrivalAirportId" example:"2"`
	DepartureAt        time.Time `json:"departureAt" example:"2026-10-15T10:00:00Z"`
	ArrivalAt          time.Time `json:"arrivalAt" example:"2026-10-15T12:00:00Z"`
	Price              int64     `json:"price" example:"5000000"`
}

type UpdateFlightRequest struct {
	AirplaneID         int64     `json:"airplaneId" example:"1"`
	DepartureAirportID int64     `json:"departureAirportId" example:"1"`
	ArrivalAirportID   int64     `json:"arrivalAirportId" example:"2"`
	DepartureAt        time.Time `json:"departureAt" example:"2026-10-15T10:00:00Z"`
	ArrivalAt          time.Time `json:"arrivalAt" example:"2026-10-15T12:00:00Z"`
	Price              int64     `json:"price" example:"5000000"`
	Status             string    `json:"status" example:"active"`
}

type FlightResponse struct {
	ID                 int64     `json:"id" example:"1"`
	AirplaneID         int64     `json:"airplaneId" example:"1"`
	DepartureAirportID int64     `json:"departureAirportId" example:"1"`
	ArrivalAirportID   int64     `json:"arrivalAirportId" example:"2"`
	DepartureAt        time.Time `json:"departureAt"`
	ArrivalAt           time.Time `json:"arrivalAt"`
	Price               int64     `json:"price" example:"5000000"`
	Capacity            int       `json:"capacity" example:"180"`
	AvailableSeats      int       `json:"availableSeats" example:"180"`
	Status              string    `json:"status" example:"active"`
}

type SearchFlightsResponse struct {
	ID                 int64     `json:"id"`
	AirplaneID         int64     `json:"airplaneId"`
	DepartureAirportID int64     `json:"departureAirportId"`
	ArrivalAirportID   int64     `json:"arrivalAirportId"`
	DepartureAt        time.Time `json:"departureAt"`
	ArrivalAt          time.Time `json:"arrivalAt"`
	Price              int64     `json:"price"`
	AvailableSeats     int       `json:"availableSeats"`
}