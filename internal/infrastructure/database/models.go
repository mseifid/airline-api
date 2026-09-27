package database

import (
	"time"
)

type AgencyModel struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	Name          string    `gorm:"column:name"`
	WalletBalance int64     `gorm:"column:wallet_balance"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (AgencyModel) TableName() string {
	return "agency"
}

type APIKeyModel struct {
	ID        int64      `gorm:"column:id;primaryKey"`
	AgencyID  int64      `gorm:"column:agency_id"`
	KeyHash   string     `gorm:"column:key_hash"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
}

func (APIKeyModel) TableName() string {
	return "api_key"
}

type AirplaneModel struct {
	ID        int16     `gorm:"column:id;primaryKey"`
	Type      string    `gorm:"column:type"`
	Capacity  int16     `gorm:"column:capacity"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AirplaneModel) TableName() string {
	return "airplane"
}