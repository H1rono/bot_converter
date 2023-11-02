package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"git.trap.jp/toki/bot_converter/model"
)

// v2 adds table converter config.
func v2() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "v2",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&model.Converter{})
		},
	}
}
