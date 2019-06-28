package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/morfien101/go-metrics-auth/config"
	"github.com/morfien101/go-metrics-auth/redisengine"
	"github.com/morfien101/go-metrics-auth/webengine"
)

var (
	// VERSION stores the version of the application
	VERSION = "0.0.2"

	flagVersion = flag.Bool("v", false, "Shows the version")
	flagHelp    = flag.Bool("h", false, "Shows the help menu")
)

func main() {
	flag.Parse()
	if *flagHelp {
		flag.PrintDefaults()
		return
	}
	if *flagVersion {
		fmt.Println(VERSION)
		return
	}

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
