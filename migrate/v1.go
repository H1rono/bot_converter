package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"git.trap.jp/toki/bot_converter/model"
)

// v1 adds table converter.
func v1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "v1",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&model.Converter{})
		},
	}
}
