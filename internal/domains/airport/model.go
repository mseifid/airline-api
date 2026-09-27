package airport

import "time"

type Airport struct {
	ID        int64
	CityID    int64
	Code      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}