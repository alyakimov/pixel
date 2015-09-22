package main


import (
    "fmt"
    "log"
    "net/http"
    "github.com/spf13/viper"
)


func main() {

    viper.SetConfigName("onepixel")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("/Users/push/Dropbox/golang/config")
    
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    router := NewRouter()

    log.Fatal(http.ListenAndServe(":8080", router))
}