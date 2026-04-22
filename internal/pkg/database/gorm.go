package database

import (
	"bookstore/internal/pkg/config"
	fileUtil "bookstore/util/file"
	"fmt"
	"path/filepath"

	"github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(injector do.Injector) (*gorm.DB, error) {
	cfg := do.MustInvoke[*config.Config](injector)
	if err := fileUtil.MkDir(filepath.Dir(cfg.Database.File)); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	db, err := gorm.Open(sqlite.Open(cfg.Database.File), &gorm.Config{})
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
