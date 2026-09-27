package database

import (
	"context"
	"fmt"

	"firefly-airline/internal/domains/agency"
	"gorm.io/gorm"
)

type AgencyRepository struct {
	db *gorm.DB
}

func NewAgencyRepository(db *gorm.DB) *AgencyRepository {
	return &AgencyRepository{
		db: db,
	}
}

var _ agency.Repository = (*AgencyRepository)(nil)

func (r *AgencyRepository) CreateWithAPIKey(ctx context.Context, a *agency.Agency, apiKey *agency.APIKey) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		agencyModel := AgencyModel{
			Name:          a.Name,
			WalletBalance: a.WalletBalance,
		}

		if err := tx.Create(&agencyModel).Error; err != nil {
			return fmt.Errorf("create agency: %w", err)
		}

		apiKeyModel := APIKeyModel{
			AgencyID: agencyModel.ID,
			KeyHash:  apiKey.KeyHash,
			RevokedAt: apiKey.RevokedAt,
		}

		if err := tx.Create(&apiKeyModel).Error; err != nil {
			return fmt.Errorf("create api key: %w", err)
		}

		a.ID = agencyModel.ID

		apiKey.ID = apiKeyModel.ID
		apiKey.AgencyID = agencyModel.ID

		return nil
	})
}