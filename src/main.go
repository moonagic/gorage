package main

import (
	"os"
	"log"
	"gorage/src/server"
	"gorage/src/config"
	"fmt"
)

const (
	version = 0.1
	buildVersion = 100
)

func main() {

	if parseArgs() {
		if err := config.LoadConfig(); err != "" {
			log.Fatal(err)
			os.Exit(0)
		}
		log.Fatal(server.StartServer(config.GetHost(), config.GetPort()))
	}

}

func parseArgs() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "help":
			printHelp()
		case "-h":
			printHelp()
		case "-help":
			printHelp()
		case "-v":
			printVersion()
		case "-V":
			printVersion()
		case "-version":
			printVersion()
		case "--version":
			printVersion()
		case "version":
			printVersion()
		case "list":
			printItemsList()
		case "delete":
			deleteTarget(os.Args[:2])
		case "start":
			return true
		}
		return false
	} else {
		return true
	}
}

func printVersion() {
	fmt.Println("gorage version" , version)
	fmt.Println("build:" , buildVersion)
}

func printHelp() {
	fmt.Println("Help info")
}

func deleteTarget(target []string) {
	fmt.Println("delete target", target)
}

func printItemsList() {
	fmt.Println("printItemsList")
}