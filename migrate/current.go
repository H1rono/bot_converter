package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"

	"git.trap.jp/toki/bot_converter/model"
)

// Migrations returns the list of all migrations.
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v1(), // v1 adds table converter.
		v2(), // v2 adds table config.
	}
}

// AllTables returns the list of all LATEST table models.
func AllTables() []interface{} {
	return []interface{}{
		&model.Converter{},
		&model.Config{},
	}
}
