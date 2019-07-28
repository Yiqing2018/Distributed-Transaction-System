package main

import (
	"log"
	"net"
	"strings"
	"time"
)

func handleListenCoordinator(recvStr string) {
	identifier := strings.Split(recvStr, " ")[0]
	switch identifier {
	case "canCommit":
		// canCommit TID
		content := strings.Split(recvStr, " ")
		tid := content[1]
		time.Sleep(1 * time.Second)
		if commitStatus[tid] {
			sendCoordinator("canCommit " + tid + " YES")
		} else {
			sendCoordinator("canCommit " + tid + " NO")
		}
	case "doCommit":
		// doCommit TID
		content := strings.Split(recvStr, " ")
		tid := content[1]
		mergeTransactions(tid)
		sendClient("COMMITOK", tidMap[tid])
		releaseLock(tid)
		log.Print("receive doCommit ", tid, " Now release Lock!")
		printLock()
		log.Print("committedMap is: ", committedMap)
	case "doAbort":
		// doAbort TID
		content := strings.Split(recvStr, " ")
		tid := content[1]
		deleteTentativeMap(tid)
		sendClient("COMMITNOTOK", tidMap[tid])
		log.Print("receive doAbort ", tid, " Now release Lock!")
		releaseLock(tid)
	default:
		log.Println("[handleListenCoordinator] wrong: **switch** receive wrong identifier")
	}
}

func handleListenClient(recvStr string, clientConn net.Conn) {
	identifier := strings.Split(recvStr, " ")[0]
	log.Println("Receive from Client: ", recvStr)
	switch identifier {
	case "BEGIN":
		// BEGIN TID
		content := strings.Split(recvStr, " ")
		tid := content[1]
		tentativeList[tid] = stringMapDeepCopy(committedMap)
		commitStatus[tid] = false
		tidMap[tid] = clientConn
	case "SET":
		// SET A.x 1 TID
		content := strings.Split(recvStr, " ")
		objectLetter := content[1][2:]
		objectValue := content[2]
		tid := content[3]

		serverletter := content[1][0:1]
		writeGetLock(tid, serverletter, objectLetter, objectValue, clientConn, true)

	case "GET":
		// GET A.x TID
		content := strings.Split(recvStr, " ")
		objectLetter := content[1][2:]
		tid := content[2]

		serverletter := content[1][0:1]
		tentativeMap := stringMapDeepCopy(tentativeList[tid])
		_, ok := tentativeMap[objectLetter]
		if !ok {
			sendClient("NOT FOUND", clientConn)
			//abort the transaction
			sendCoordinator("ABORT" + " " + tid)
			//release the lock
		} else {
			readGetLock(tid, serverletter, objectLetter, clientConn, true)
		}

	case "COMMIT":
		// COMMIT TID
		content := strings.Split(recvStr, " ")
		tid := content[1]
		commitStatus[tid] = true

	default:
		log.Println("[handleListenClient] wrong: **switch** receive wrong identifier")
	}

	printLock()

}

func mergeTransactions(tid string) {
	// merge the committedMap with tentativeList[tid]
	tentativeMap := stringMapDeepCopy(tentativeList[tid])
	for objectLetter, objectValue := range tentativeMap {
		committedMap[objectLetter] = objectValue
	}
}

func deleteTentativeMap(tid string) {
	delete(tentativeList, tid)
}

func stringMapDeepCopy(map1 map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range map1 {
		newMap[k] = v
	}

	return newMap
}

func printLock() {
	log.Print("----------------LOCK-----------------------")
	log.Print("lock information: ", serverLock.lockInfo)
	log.Print("holder information: ", serverLock.holder)
	log.Print("---------------------------------------")
}
