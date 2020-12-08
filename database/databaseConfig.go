package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	DB *gorm.DB
}

// DBInit is open a database and save the reference to `Config` struct.
func (d *Config) DBInit() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPass))

	if err != nil {
		log.Fatal("db err: ", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.LogMode(true)
	d.DB = db
}

// TestDBInit will create a temporarily database for running testing cases
func (d *Config) TestDBInit() {
	//"sqlite3", "file::memory:?mode=memory&cache=shared"
	// testDb, err := gorm.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	testDb, err := gorm.Open("sqlite3", "sqlite-database.db")
	if err != nil {
		log.Fatal("db err: ", err)
	}

	testDb.DB().SetMaxIdleConns(10)
	testDb.DB().SetMaxOpenConns(100)
	testDb.DB().SetConnMaxLifetime(time.Hour)
	testDb.LogMode(true)
	d.DB = testDb
}

// TestDBFree delete the database after running testing cases.
func (d *Config) TestDBFree(testDb *gorm.DB) error {
	if err := testDb.Close(); err != nil {
		return err
	}
	err := os.Remove("sqlite-database.db")
	return err
}

// DBClose ...
func (d *Config) DBClose(db *gorm.DB) {
	if err := db.Close(); err != nil {
		log.Error(err)
	}
}
