package airplane

import "time"

type Airplane struct {
	ID        int64
	Type      string
	Capacity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}