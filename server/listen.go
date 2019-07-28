package main

import (
	"log"
	"net"
)

func listenCoordinator(listenConn net.Conn) {
	for {
		buf := make([]byte, 1024)
		length, err := listenConn.Read(buf)
		if err != nil {
			log.Println("[listenCoordinator] wrong: ", err)
			break
		}
		recvStr := string(buf[0:length])
		go handleListenCoordinator(recvStr)
	} // infinite loop
}

func pickClient() {
	listenPort := ":" + serverPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenPort)
	if err != nil {
		log.Println("[pickClient] wrong: ", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("[pickClient] wrong: ", err)
		return
	}

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go listenClient(clientConn)
	} // infinite loop
}

func listenClient(clientConn net.Conn) {
	for {
		buf := make([]byte, 1024)
		length, err := clientConn.Read(buf)
		if err != nil {
			log.Println("[listenClient] wrong: ", err)
			break
		}
		recvStr := string(buf[0:length])
		go handleListenClient(recvStr, clientConn)
	} // infinite loop
}
