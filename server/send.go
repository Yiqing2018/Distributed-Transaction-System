package main

import (
	"log"
	"net"
)

func sendClient(message string, clientConn net.Conn) {
	for {
		_, err := clientConn.Write([]byte(message))
		if err != nil {
			log.Println("[sendClient] wrong: ", err)
			continue
		}
		break
	} // infinite loop with break
}

func sendCoordinator(message string) {
	for {
		_, err := coordinatorConn.Write([]byte(message))
		if err != nil {
			log.Println("[sendCoordinator] wrong: ", err)
			continue
		}
		break
	} // infinite loop with break
}
