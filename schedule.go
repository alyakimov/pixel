package main

import (
	"fmt"
	"log"
	"sync"
)

func campaingsCacheUpdater() {
	log.Println("campaingsCacheUpdater running")

	var mutex = &sync.Mutex{}

	db := GetConnection()
	temp_campaings, err := GetAllCampaign(db)

	if err != nil {
		panic(fmt.Errorf("Fatal error running update cache campaings\n", err))
	} else {
		mutex.Lock()
		campaings = temp_campaings
		log.Println("campaingsCacheUpdater success")
		mutex.Unlock()
	}
}
