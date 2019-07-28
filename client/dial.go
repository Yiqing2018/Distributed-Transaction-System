package main

import (
	"log"
	"net"
	"sync"
)

func dialServer() {
	// dial to five servers
	var wg sync.WaitGroup

	for idx := range serverAddrs {
		dialAddress := serverAddrs[idx]
		for {
			serverConn, err := net.Dial("tcp", dialAddress)
			if err != nil {
				log.Println("[dialServer] wrong: ", err)
				continue
			}
			serverLetter := string(idx + 65)
			serverConnMap[serverLetter] = serverConn
			// go listenServer(serverConn)

			wg.Add(1)
			go func() {
				defer wg.Done()
				listenServer(serverConn)
			}()

			break
		} // infinite loop with break

	}

	wg.Wait()

}

func dialCoordinator() {
	dialAddress := coordinatorAddress

	for {
		coordinatorCon, err := net.Dial("tcp", dialAddress)
		if err != nil {
			log.Println("[dialCoordinator] wrong: ", err)
			continue
		}
		coordinatorConn = coordinatorCon
		break
	} // infinite loop with break

	listenCoordinator()

}
