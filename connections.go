package main


import (
    "log"
    "database/sql"
    "github.com/spf13/viper"
    _ "github.com/go-sql-driver/mysql"
)


var db *sql.DB = nil


func GetConnection() *sql.DB {

    if db != nil {
        return db
    
    } else {
        db, err := sql.Open("mysql", viper.GetString("connections.onepixel.dsl"))

        if err != nil {
            log.Fatalf("Error on initializing database connection: %s", err.Error())
        }

        db.SetMaxIdleConns(viper.GetInt("connections.onepixel.maxIdleConnection"))

        err = db.Ping()
        if err != nil {
            log.Fatalf("Error on opening database connection: %s", err.Error())
        }

        return db
    }
}