package model

import (
	"github.com/gofrs/uuid"
)

type Config struct {
	ConverterID        uuid.UUID   `gorm:"type:char(36);not null;primaryKey"`
	PushBranchFilter   StringSlice `gorm:"type:text;not null"`
	PREventTypesFilter StringSlice `gorm:"type:text;not null"`
}
