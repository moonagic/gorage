package main

import (
	"imagesStorage/src/server"
	"log"
	"imagesStorage/src/config"
	"os"
)

func main() {

	if err := config.LoadConfig(); err != "" {
		log.Fatal(err)
		os.Exit(0)
	}
	log.Fatal(server.StartServer(config.GetHost(), config.GetPort()))
}

