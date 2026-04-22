package database

import (
	"bookstore/internal/pkg/config"

	"github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(injector do.Injector) (*gorm.DB, error) {
	cfg := do.MustInvoke[*config.Config](injector)
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}
