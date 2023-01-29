package models

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` //add soft delete in gorm
}

type BinaryBase struct {
	ID        BINARY16       `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` //add soft delete in gorm
}

type BaseModel struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	CreateDateTime time.Time      `gorm:"autoCreateTime;column:CreateDateTime" json:"created_datetime"`
	UpdateDateTime time.Time      `gorm:"autoUpdateTime;column:UpdateDateTime" json:"updated_datetime"`
	DeleteDateTime gorm.DeletedAt `gorm:"column:DeleteDateTime" json:"deleted_datetime"`
	DeleteFlg      bool           `gorm:"column:DeleteFlg;softDelete:flag" json:"is_deleted"`
}
