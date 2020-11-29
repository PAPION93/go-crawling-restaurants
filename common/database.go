package common

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/viper"
)

type Database struct {
	DB *sql.DB
}

func Init() *sql.DB {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}
