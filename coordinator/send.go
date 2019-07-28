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
		log.Println("[sendClient] Send msg to client: ", message)
		break
	} // infinite loop with break
}

func sendServer(message string) {
	for serverConn := range serverConnSet {
		for {
			_, err := serverConn.Write([]byte(message))
			if err != nil {
				log.Println("[sendServer] wrong: ", err)
				continue
			}
			break
		} // infinite loop with break
	}
}
