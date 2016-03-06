package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {

	viper.SetConfigName("onepixel")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(viper.GetString("port"), router))
}
