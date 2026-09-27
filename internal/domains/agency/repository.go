package agency

import "context"

type Repository interface {
	CreateWithAPIKey(ctx context.Context, agency *Agency, apiKey *APIKey) error
	GetAgencyIDByAPIKeyHash(ctx context.Context, keyHash string) (int64, error)
	AddBalance(ctx context.Context, agencyID int64,amount int64) (int64, error)
}