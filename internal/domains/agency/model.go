package agency

import "time"

type Agency struct {
	ID            int64
	Name          string
	WalletBalance int64
}

type APIKey struct {
	ID        int64
	AgencyID  int64
	KeyHash   string
	CreatedAt time.Time
	RevokedAt *time.Time
}
