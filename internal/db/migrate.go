package db

import (
	"database/sql"
	"log"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&interfaces.User{},
	)
	if err != nil {
		return err
	}

	return nil
}

func RawDB(db *gorm.DB) *sql.DB {
	rawDB, err := db.DB()
	if err != nil {
		log.Panicf("Unable to get raw sql.DB %s\n", err.Error())
	}

	return rawDB
}
