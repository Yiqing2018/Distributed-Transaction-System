package main

import (
	"log"
	"strings"
)

func sendCoordinator(message string) {
	for {
		_, err := coordinatorConn.Write([]byte(message))
		if err != nil {
			log.Println("[sendCoordinator] wrong: ", err)
			continue
		}
		log.Println("[sendCoordinator] Send msg to coor:", message)
		break
	} // infinite loop with break
}

func sendServer(message string) {
	identifier := strings.Split(message, " ")[0]
	if identifier == "SET" || identifier == "GET" {
		log.Println("send to server: ", message)
		content := strings.Split(message, " ")
		serverLetter := content[1][0:1]
		serverConn := serverConnMap[serverLetter]
		for {
			_, err := serverConn.Write([]byte(message))
			if err != nil {
				log.Println("[sendServer] wrong: ", err)
				continue
			}
			break
		} // infinite loop with break

	} else if identifier == "BEGIN" || identifier == "COMMIT" || identifier == "ABORT" {
		for _, serverConn := range serverConnMap {
			log.Println("send to server: ", message)
			for {
				_, err := serverConn.Write([]byte(message))
				if err != nil {
					log.Println("[sendServer] wrong: ", err)
					continue
				}
				break
			} // infinite loop with break
		}
	} else {
		log.Println("[sendServer] wrong: **if** is wrong.")
	}

}
