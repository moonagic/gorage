package main

import (
	"imagesStorage/src/server"
	"log"
)

func main() {
	if error := server.StartServer("", ""); error != nil {
		log.Fatal(error)
	}
}

