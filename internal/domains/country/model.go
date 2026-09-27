package country

import "time"

type Country struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}