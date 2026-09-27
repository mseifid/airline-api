package flight

import "errors"

var (
	ErrNotFound = errors.New("flight not found")
	ErrInvalidAirports = errors.New(
		"departure and arrival airports must be different",
	)
	ErrInvalidSchedule = errors.New(
		"arrival time must be after departure time",
	)
	ErrInvalidPrice = errors.New(
		"price cannot be negative",
	)
)