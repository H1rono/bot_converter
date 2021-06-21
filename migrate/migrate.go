package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrate executes the migration.
func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              191,
		UseTransaction:            false,
		ValidateUnknownMigrations: true,
	}, Migrations())

	m.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(AllTables()...); err != nil {
			return err
		}
		return nil
	})

	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
