package main

import (
	"log"
	"net"
	"time"
)

func pickClientAndServer() {
	listenPort := ":" + coordinatorPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenPort)
	if err != nil {
		log.Println("[pickClientAndServer] wrong: ", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("[pickClientAndServer] wrong: ", err)
		return
	}
	addressCount := 0

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		if addressCount < 5 {
			serverConnSet[conn] = false
			addressCount = addressCount + 1
			if addressCount == 5 {
				log.Println("[pickClientAndServer] Accept 5 servers.")
			}
		}

		//开始监听这个conn 无论是和server的还是和client的
		go listenAll(conn)
	} // infinite loop
}

func listenAll(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		length, err := conn.Read(buf)
		if err != nil {
			log.Print("conn: ", conn)
			log.Println("[listenAll] wrong: ", err)
			time.Sleep(1 * time.Second)
			break
		}
		recvStr := string(buf[0:length])
		log.Println("received message:", recvStr)
		handleLlistenAll(recvStr, conn)
	} // infinite loop
}
