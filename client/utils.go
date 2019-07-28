package main

import (
	"fmt"
	"log"
	"strings"
)

func setTID(tid string) {
	TID = tid
}

func printUserResponse(response string) {
	fmt.Println(response)
}

func handleListenCoordinator(recvStr string) {
	splits := strings.Split(recvStr, " ")
	if splits[0] == "TID" {
		// listen == TID 123
		setTID(splits[1])
		log.Println("Set TID successfully!")
		sendServer("BEGIN" + " " + TID)
		log.Println("Send msg to Server: ", "BEGIN"+" "+TID)
		printUserResponse("OK")
	} else if splits[0] == "COMMIT" {
		// listen == COMMIT TID
		printUserResponse("COMMIT OK")
	} else if splits[0] == "ABORTED" {
		// listen == ABORTED TID
		printUserResponse("ABORTED")
	} else {
		log.Println("[handleListenCoordinator] wrong: **if** is wrong", splits[0])
	}
}

func handleListenServer(recvStr string) {
	splits := strings.Split(recvStr, " ")
	if splits[0] == "SET" {
		// listen == SET A.x 1 OK
		printUserResponse("OK")
	} else if splits[0] == "GET" {
		// listen == GET A.x 1
		printUserResponse(splits[1] + " " + "=" + " " + splits[2])
	} else if splits[0] == "NOT" {
		// listen == NOT FOUND
		printUserResponse("NOT FOUND")
	} else if splits[0] == "COMMITOK" {
		// listen == COMMITOK
		if beginFlag {
			beginFlag = false
			printUserResponse("COMMIT OK")
		}
	} else if splits[0] == "COMMITNOTOK" {
		// listen == COMMITNOTOK
		if beginFlag {
			beginFlag = false
			printUserResponse("ABORTED")
		}
	} else {
		log.Println("[handleListenServer] wrong: **if** is wrong", splits[0])
	}
}

func handleListenUser(input string) {
	identifier := strings.Split(input, " ")[0]
	switch identifier {
	case "BEGIN":
		// input == BEGIN
		// send = "BEGIN"
		sendCoordinator("BEGIN")
		beginFlag = true
		//time.Sleep(3 * time.Second)
	case "SET":
		// input == SET A.x 1
		// send = "SET A.x 1 TID"
		sendServer(input + " " + TID)
	case "GET":
		// input == GET A.x
		// send = "GET A.x TID"
		sendServer(input + " " + TID)
	case "COMMIT":
		// input == COMMIT
		// send = "COMMIT TID"
		sendServer("COMMIT" + " " + TID)
		sendCoordinator("COMMIT" + " " + TID)
	case "ABORT":
		// input == ABORT
		// send = "ABORT TID"
		// sendServer("ABORT" + " " + TID)
		sendCoordinator("ABORT" + " " + TID)
	default:
		log.Println("[handleListenUser] wrong: **switch** receive wrong identifier")
	}
}
