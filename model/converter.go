package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Converter struct {
	ID        uuid.UUID      `gorm:"type:char(36);not null;primaryKey"`
	CreatorID uuid.UUID      `gorm:"type:char(36);not null"`
	ChannelID uuid.UUID      `gorm:"type:char(36);not null"`
	Secret    sql.NullString `gorm:"type:text"`
	CreatedAt time.Time      `gorm:"precision:6;not null;default:NOW()"`
	DeletedAt gorm.DeletedAt `gorm:"precision:6"`
}
