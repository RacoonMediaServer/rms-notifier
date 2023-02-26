package db

import (
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

func Connect(config config.Database) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(config.Path))
	if err != nil {
		return nil, err
	}
	return &Database{conn: db}, nil
}
