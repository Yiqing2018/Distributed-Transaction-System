package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func handleLlistenAll(recvStr string, conn net.Conn) {
	identifier := strings.Split(recvStr, " ")[0]
	switch identifier {
	case "BEGIN":
		// BEGIN (from client)
		log.Println("[handleLlistenAll] Enter Switch: BEGIN")
		tid := generateTID()
		canCommitCount[tid] = 0
		canAbortCount[tid] = false
		sendClient("TID "+tid, conn)
	case "COMMIT":
		// COMMIT TID (from client)
		content := strings.Split(recvStr, " ")
		tid := content[1]
		sendServer("canCommit" + " " + tid)
	case "ABORT":
		// ABORT TID (from client or server)
		content := strings.Split(recvStr, " ")
		tid := content[1]
		graphForever.handleDelete(tid)
		sendServer("doAbort" + " " + tid)
	case "canCommit":
		// canCommit TID YES (from server)
		// canCommit TID NO (from server)
		content := strings.Split(recvStr, " ")
		tid := content[1]
		action := content[2]
		if action == "YES" {
			canCommitCount[tid] = canCommitCount[tid] + 1
			if canCommitCount[tid] == 5 {
				graphForever.handleDelete(tid)
				sendServer("doCommit" + " " + tid)
			}
		} else {
			if !canAbortCount[tid] {
				canAbortCount[tid] = true
				log.Println("Receive NO, and send doAbort to 5 servers")
				graphForever.handleDelete(tid)
				sendServer("doAbort" + " " + tid)
			}
		}
	case "graph":
		// graph start end (from server)
		content := strings.Split(recvStr, " ")
		testEdge := content[1] + " " + content[2]
		graphForever.addTestEdge(testEdge)
		graphForever.checkCircle()
	case "release":
		// release tid
		content := strings.Split(recvStr, " ")
		releaseTID := content[1]
		graphForever.handleDelete(releaseTID)
	default:
		log.Println("[handleLlistenAll] wrong: **switch** receive wrong identifier ", identifier)
	}
}

// no use
func generateTIDRandom() string {
	returnRandom := "default"
	for {
		rand.Seed(int64(time.Now().UnixNano()))
		randomNumber := rand.Intn(10000)
		randomString := strconv.Itoa(randomNumber)
		_, ok := tidSet[randomString]
		if ok {
			continue
		} else {
			returnRandom = randomString
			break
		}
	} // infinite loop with break

	tidSet[returnRandom] = false
	return returnRandom
}

func generateTID() string {
	randomString := strconv.Itoa(tidNO)
	tidNO = tidNO + 1
	return randomString
}
