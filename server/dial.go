package main

import (
	"log"
	"net"
	"sync"
)

func dialCoordinator() {
	dialAddress := coordinatorAddress
	var wg sync.WaitGroup
	for {
		coordinatorCon, err := net.Dial("tcp", dialAddress)
		if err != nil {
			log.Println("[dialCoordinator] wrong: ", err)
			continue
		}
		coordinatorConn = coordinatorCon
		// go listenCoordinator(coordinatorConn)

		wg.Add(1)
		go func() {
			defer wg.Done()
			listenCoordinator(coordinatorConn)
		}()

		break
	} // infinite loop with break
	wg.Wait()

}
