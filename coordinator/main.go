package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	graphForever.initGraph()
	coordinatorPort = "6666"
	file, err := os.Create("coordinator.log")
	defer file.Close()
	if err != nil {
		fmt.Println("[main] wrong: ", err)
		return
	}
	log.SetOutput(file)
	defer file.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pickClientAndServer()
	}()
	wg.Wait()

}
