package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

var db *sql.DB = nil

func GetConnection() *sql.DB {

	if db != nil {
		return db

	} else {

		var err error
		db, err = sql.Open("mysql", viper.GetString("connections.onepixel.dsl"))

		if err != nil {
			log.Fatalf("Error on initializing database connection: %s", err.Error())
		}

		db.SetMaxIdleConns(viper.GetInt("connections.onepixel.maxIdleConnection"))
		db.SetMaxOpenConns(viper.GetInt("connections.onepixel.maxOpenConnection"))

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error on opening database connection: %s", err.Error())
		}

		return db
	}
}
