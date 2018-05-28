package main

import (
	"gorage/src/server"
	"log"
	"gorage/src/config"
	"os"
)

func main() {

	if err := config.LoadConfig(); err != "" {
		log.Fatal(err)
		os.Exit(0)
	}
	log.Fatal(server.StartServer(config.GetHost(), config.GetPort()))
}

