package province

import "time"

type Province struct {
	ID        int64
	CountryID int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}