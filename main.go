package main

import (
	"fmt"
	"log"

	"github.com/morfien101/go-metrics-auth/config"
	"github.com/morfien101/go-metrics-auth/redisengine"
	"github.com/morfien101/go-metrics-auth/webengine"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("Failed to make config. Error:", err)
	}
	re := redisengine.New(config.Redis)
	if err := re.Start(); err != nil {
		log.Fatal("Failed to create redis Engine. Error:", err)
	}

	webengine := webengine.New(config.WebServer, re)
	log.Println("Starting Web server")
	err = <-webengine.Start()
	if err != nil {
		fmt.Println(err)
	}
}
