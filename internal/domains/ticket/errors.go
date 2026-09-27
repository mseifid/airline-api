package ticket

import "errors"

var (
	ErrInvalidPassengerCount = errors.New("at least one passenger is required")
	ErrFlightNotFound        = errors.New("flight not found")
	ErrNotEnoughSeats        = errors.New("not enough available seats")
	ErrInsufficientBalance   = errors.New("insufficient wallet balance")
	ErrTicketNotFound        = errors.New("ticket not found")
	ErrTicketAlreadyCancelled = errors.New("ticket is already cancelled")
	ErrFlightAlreadyDeparted = errors.New("flight has already departed")
)