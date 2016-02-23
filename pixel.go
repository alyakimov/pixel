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
    viper.AddConfigPath(".")
    
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    } 

    uri := viper.GetString("amqp.uri")
    exchange := viper.GetString("amqp.exchange")
    exchangeType := viper.GetString("amqp.exchange_type")
    queue := viper.GetString("amqp.queue")
    bindingKey := viper.GetString("amqp.binding_key")
    consumerTag := viper.GetString("amqp.consumer_tag")

    _, err = NewConsumer(uri, exchange, exchangeType, queue, bindingKey, consumerTag)

    if err != nil {
        log.Fatalf("%s", err)
    }    

    router := NewRouter()

    log.Fatal(http.ListenAndServe(viper.GetString("port"), router))
}