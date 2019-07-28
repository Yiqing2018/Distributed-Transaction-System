package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Create("client" + ".log")
	defer file.Close()
	if err != nil {
		fmt.Println("[main] wrong: ", err)
		return
	}
	log.SetOutput(file)
	defer file.Close()

	go dialServer()
	go dialCoordinator()
	listenUser()
}
