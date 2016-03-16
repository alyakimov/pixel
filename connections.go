package main

import (
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

var db *mgo.Session = nil

func GetConnection() *mgo.Session {

	if db != nil {
		return db

	} else {

		var err error

		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:   []string{viper.GetString("connections.mongo.hosts")},
			Timeout: 60 * time.Second,
			Database: viper.GetString("connections.mongo.auth_database"),
			Username: viper.GetString("connections.mongo.auth_user_name"),
			Password: viper.GetString("connections.mongo.auth_password"),
		}

		db, err = mgo.DialWithInfo(mongoDBDialInfo)

		if err != nil {
			log.Fatalf("Error on initializing database connection: %s", err.Error())
		}

		db.SetMode(mgo.Monotonic, true)

		return db
	}
}
