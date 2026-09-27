package ticket

import (
	"context"
	"time"
)

type Repository interface {
	Purchase(
		ctx context.Context,
		agencyID int64,
		flightID int64,
		passengers []Passenger,
	) (*Ticket, error)

	ListByFlightAndAgency(
		ctx context.Context,
		flightID int64,
		agencyID int64,
	) ([]Ticket, error)

	GetForCancellation(
		ctx context.Context,
		ticketID int64,
		agencyID int64,
	) (*Ticket, error)

	Cancel(
		ctx context.Context,
		ticketID int64,
		agencyID int64,
		refundAmount int64,
		now time.Time,
	) (*Ticket, error)

	GetByID(
		ctx context.Context,
		ticketID int64,
		agencyID int64,
	) (*Ticket, error)
}
