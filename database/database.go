package database

import (
	"go-orm-vs-sql/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	// GORM用データベース接続
	gormDB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// GORMの自動マイグレーション
	if err := gormDB.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	return gormDB, nil
}
