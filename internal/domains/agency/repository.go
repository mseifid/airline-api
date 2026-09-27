package agency

import "context"

type Repository interface {
	CreateWithAPIKey(ctx context.Context, agency *Agency, apiKey *APIKey) error
	GetAgencyIDByAPIKeyHash(ctx context.Context, keyHash string) (int64, error)
}