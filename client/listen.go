package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func listenServer(listenConn net.Conn) {
	for {
		buf := make([]byte, 1024)
		length, err := listenConn.Read(buf)
		if err != nil {
			log.Println("[listenServer] wrong: ", err)
			break
		}
		recvStr := string(buf[0:length])
		go handleListenServer(recvStr)
	} // infinite loop
}

func listenCoordinator() {
	for {
		buf := make([]byte, 1024)
		length, err := coordinatorConn.Read(buf)
		if err != nil {
			log.Println("[listenCoordinator] wrong: ", err)
			break
		}
		recvStr := string(buf[0:length])
		log.Println("[listenCoordinator] Get msg from Coordinator: ", recvStr)

		handleListenCoordinator(recvStr)

	} // infinite loop

}

func listenUser() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		if input != "" {
			handleListenUser(input)
		}
		reader.Reset(os.Stdin)
	} // infinite loop
}
