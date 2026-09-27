package agency

import "context"

type Repository interface {
	CreateWithAPIKey(ctx context.Context, agency *Agency, apiKey *APIKey) error
}