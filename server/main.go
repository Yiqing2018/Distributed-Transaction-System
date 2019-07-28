package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	file, err := os.Create("server.log")
	if err != nil {
		fmt.Println("fail to create log file")
		return
	}
	log.SetOutput(file)
	defer file.Close()

	go func() {
		defer wg.Done()
		dialCoordinator()
	}()
	pickClient()

	wg.Wait()

}
