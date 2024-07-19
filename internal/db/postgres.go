package db

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	// new db
	db := connectToPostgres()

	if gin.Mode() == gin.ReleaseMode {
		db.Logger.LogMode(0)
	}

	rawDB := RawDB(db)

	rawDB.SetMaxIdleConns(20)
	rawDB.SetMaxOpenConns(100)

	// migrate models
	err := Migrate(db)
	if err != nil {
		log.Panicf("Unable to migrate models %s\n", err.Error())
	}
	log.Println("Successfully Migrated Models.")

	return db
}

func connectToPostgres() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")

	var counts int

	for {
		connection, err := openPostgres(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Fatal(err)
		}

		log.Println("Backing off for three seconds....")
		time.Sleep(3 * time.Second)
		continue
	}
}

func openPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// return *sql.DB from db(*gorm.DB) to enable Ping()
	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// ping database
	err = gormDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
