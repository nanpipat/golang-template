package migrations

import (
	"github.com/nanpipat/golang-template-hexagonal/internal/core/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	tables := []interface{}{
		&domain.User{},
	}

	err := db.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	return nil
}
