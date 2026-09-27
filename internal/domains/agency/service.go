package agency

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	name string,
	walletBalance int64,
) (Agency, string, error) {
	if name == "" {
		return Agency{}, "", fmt.Errorf("name is required")
	}

	if walletBalance < 0 {
		return Agency{}, "", fmt.Errorf("wallet balance cannot be negative")
	}

	rawKey, err := generateAPIKey()
	if err != nil {
		return Agency{}, "", fmt.Errorf("generate api key: %w", err)
	}

	hash := sha256.Sum256([]byte(rawKey))

	agency := Agency{
		Name:          name,
		WalletBalance: walletBalance,
	}

	apiKey := APIKey{
		KeyHash: hex.EncodeToString(hash[:]),
	}

	if err := s.repository.CreateWithAPIKey(
		ctx,
		&agency,
		&apiKey,
	); err != nil {
		return Agency{}, "", fmt.Errorf("create agency: %w", err)
	}

	return agency, rawKey, nil
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "ak_" + hex.EncodeToString(b), nil
}