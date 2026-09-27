package agency

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func (s *Service) ValidateAPIKey(ctx context.Context, rawKey string) (int64, error) {
	if rawKey == "" {
		return 0, errors.New("api key is required")
	}

	hash := sha256.Sum256([]byte(rawKey))

	keyHash := hex.EncodeToString(hash[:])

	agencyID, err := s.repository.GetAgencyIDByAPIKeyHash(
		ctx,
		keyHash,
	)
	if err != nil {
		return 0, err
	}

	return agencyID, nil
}

func (s *Service) AddBalance(ctx context.Context,agencyID int64,amount int64) (int64, error) {
	if agencyID <= 0 {
		return 0, errors.New("invalid agency")
	}

	if amount <= 0 {
		return 0, ErrInvalidDepositAmount
	}

	return s.repository.AddBalance(ctx, agencyID, amount)
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "ak_" + hex.EncodeToString(b), nil
}