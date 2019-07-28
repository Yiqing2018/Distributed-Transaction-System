package main

import (
	"log"
	"net"
	"strings"
)

//GET
func readGetLock(tid string, serverLetter string, objectLetter string, clientConn net.Conn, putIntoBuffer bool) bool {

	canDoOperation := false
	serverLock.Lock()
	_, ok := serverLock.lockInfo[objectLetter]
	if !ok {
		//no one is holding the lock
		serverLock.lockInfo[objectLetter] = 0
		var mySlice []string
		mySlice = append(mySlice, tid)
		serverLock.holder[objectLetter] = mySlice
		//do the operation!!!
		canDoOperation = true

	} else {
		//someone is holding the lock ---  maybe myself
		lockType := serverLock.lockInfo[objectLetter]
		if lockType == 0 {
			mySlice := serverLock.holder[objectLetter]
			//check if I already got the readLock
			flag := false
			for _, holder := range mySlice {
				if holder == tid {
					flag = true
				}
			}
			if !flag {
				mySlice = append(mySlice, tid)
				serverLock.holder[objectLetter] = mySlice
			}
			//do the operation!!!
			canDoOperation = true

		} else {
			//someone is holding the write lock ---  maybe myself
			mySlice := serverLock.holder[objectLetter]
			//check if I already got the writeLock
			holder := mySlice[0]
			if holder == tid {
				//do the operation!!!
				canDoOperation = true
			} else {
				//someone else is holding the writeLock, I can do nothing... waiting!
				// send "wait-for" information to coordinator
				if putIntoBuffer {
					pushToQueue(tid+" "+serverLetter+" "+objectLetter, clientConn)
				}
				tobeSend := "graph " + holder + " " + tid
				log.Print("waitfor info ", tobeSend)
				sendCoordinator(tobeSend)
			}
		}
	}
	serverLock.Unlock()

	if canDoOperation {
		objectVal := tentativeList[tid][objectLetter]
		SendToClient := "GET " + serverLetter + "." + objectLetter + " " + objectVal
		sendClient(SendToClient, clientConn)
	}
	return canDoOperation

}

//SET
func writeGetLock(tid string, serverLetter string, objectLetter string, objectValue string, clientConn net.Conn, putIntoBuffer bool) bool {
	//log.Print("write0")
	canDoOperation := false
	serverLock.Lock()
	_, ok := serverLock.lockInfo[objectLetter]
	if !ok {
		//log.Print("write1")
		//no one is holding the lock
		serverLock.lockInfo[objectLetter] = 1
		var mySlice []string
		mySlice = append(mySlice, tid)
		serverLock.holder[objectLetter] = mySlice

		//do the operation!!!
		canDoOperation = true

	} else {
		//log.Print("write2")
		//someone is holding the lock ---  maybe myself
		lockType := serverLock.lockInfo[objectLetter]
		if lockType == 0 {
			// log.Print("write3")
			mySlice := serverLock.holder[objectLetter]
			flag := false
			for _, holder := range mySlice {

				if holder != tid {
					tobeSend := "graph " + holder + " " + tid
					log.Print("waitfor info ", tobeSend)
					sendCoordinator(tobeSend)
					flag = true
				}
			}
			if flag {
				//someone else is holding the readLock, I can do nothing... waiting!
				if putIntoBuffer {
					pushToQueue(tid+" "+serverLetter+" "+objectLetter+" "+objectValue, clientConn)
				}

			} else {
				// log.Print("write4")
				// i am the only one who is holding the readLock, upgrade the lock!
				serverLock.lockInfo[objectLetter] = 1
				//do the operation
				canDoOperation = true

			}
		} else {
			// log.Print("write5")
			//write lock, only one holder
			mySlice := serverLock.holder[objectLetter]
			holder := mySlice[0]
			if holder == tid {
				//do the operation

				canDoOperation = true
			} else {
				//someone else is holding the writeLock, I can do nothing... waiting!
				// log.Print("write6")
				if putIntoBuffer {
					pushToQueue(tid+" "+serverLetter+" "+objectLetter+" "+objectValue, clientConn)
				}
				// log.Print("write7")
				tobeSend := "graph " + holder + " " + tid
				// log.Print("write8")
				sendCoordinator(tobeSend)
				// log.Print("write9")
				log.Print("waitfor info ", tobeSend)
				// log.Print("write10")

			}

		}

	}
	serverLock.Unlock()
	if canDoOperation {
		tentativeList[tid][objectLetter] = objectValue
		SendToClient := "SET " + serverLetter + "." + objectLetter + " " + objectValue + " OK"
		sendClient(SendToClient, clientConn)
	}
	return canDoOperation
}

//realseLock
// func releaseLock(tid string) {
// 	serverLock.Lock()
// 	for object, holderSlice := range serverLock.holder {
// 		//how many holders
// 		count := len(holderSlice)
// 		//go through the holderList
// 		for idx, tx := range holderSlice {
// 			if tx == tid {
// 				//I am the only one holding the lock
// 				if count == 1 {
// 					delete(serverLock.lockInfo, object)
// 					delete(serverLock.holder, object)
// 				} else {
// 					//I am holding the lock with someone else
// 					newSlice := append(holderSlice[:idx], holderSlice[idx+1:]...)
// 					serverLock.holder[object] = newSlice
// 				}
// 			}
// 		}

// 	}
// 	serverLock.Unlock()

// 	buffer.Lock()
// 	for idx, opt := range buffer.operations {
// 		tidInBuffer := strings.Split(opt, "")[0]
// 		if tidInBuffer == tid {
// 			buffer.operations = append(buffer.operations[:idx], buffer.operations[idx+1:]...)
// 			buffer.conns = append(buffer.conns[:idx], buffer.conns[idx+1:]...)
// 		}
// 	}
// 	buffer.Unlock()

// 	checkBuffer()
// }

func releaseLock(tid string) {
	//log.Print("releaseLock 0")
	serverLock.Lock()
	for object, holderSlice := range serverLock.holder {
		//log.Print("releaseLock 1")
		//how many holders
		count := len(holderSlice)
		var tempSlice []string
		//go through the holderList
		//log.Print("releaseLock 2")
		for _, tx := range holderSlice {

			if tx == tid {
				//I am the only one holding the lock
				if count == 1 {
					delete(serverLock.lockInfo, object)
					delete(serverLock.holder, object)
				} else {
					//I am holding the lock with someone else
					//sendCoordinator("release " + tid)
					continue
				}
			}
			tempSlice = append(tempSlice, tx)
		}
		if count != 1 {
			serverLock.holder[object] = tempSlice
		}
		//log.Print("releaseLock 3")
	}

	//log.Print("releaseLock 4")
	serverLock.Unlock()
	//log.Print("releaseLock JOURNEY MARK")
	buffer.Lock()
	//log.Print("releaseLock 5")

	var tempOperations []string
	var tempConns []net.Conn
	for idx, opt := range buffer.operations {
		tidInBuffer := strings.Split(opt, " ")[0]
		if tidInBuffer == tid {
			continue
		}
		tempOperations = append(tempOperations, buffer.operations[idx])
		tempConns = append(tempConns, buffer.conns[idx])

	}
	buffer.operations = tempOperations
	buffer.conns = tempConns

	//log.Print("releaseLock 6")
	buffer.Unlock()
	//log.Print("releaseLock 7")

	checkBuffer()
	//log.Print("releaseLock 8")
}

//when we releaseLock, check the buffer
func checkBuffer() {
	//log.Print("checkBuffer 1")
	buffer.Lock()
	//log.Print("checkBuffer 2")
	length := len(buffer.operations)
	if length == 0 {
		buffer.Unlock()
		return
	}
	//log.Print("checkBuffer 3")
	for idx, opt := range buffer.operations {
		clientConn := buffer.conns[idx]
		log.Print("buffer has operation: ", opt)
		arrs := strings.Split(opt, " ")
		//read operation
		if len(arrs) == 3 {
			tid := arrs[0]
			serverLetter := arrs[1]
			objectLetter := arrs[2]
			succss := readGetLock(tid, serverLetter, objectLetter, clientConn, false)
			if succss {
				//remove from the buffer?????
				buffer.operations = append(buffer.operations[:idx], buffer.operations[idx+1:]...)
				buffer.conns = append(buffer.conns[:idx], buffer.conns[idx+1:]...)
			}

			//write operation
		} else {
			tid := arrs[0]
			serverLetter := arrs[1]
			objectLetter := arrs[2]
			objectValue := arrs[3]
			succss := writeGetLock(tid, serverLetter, objectLetter, objectValue, clientConn, false)
			if succss {
				//remove from the buffer
				buffer.operations = append(buffer.operations[:idx], buffer.operations[idx+1:]...)
				buffer.conns = append(buffer.conns[:idx], buffer.conns[idx+1:]...)
			}
		}

	}

	//log.Print("checkBuffer 4")
	buffer.Unlock()
	//log.Print("checkBuffer 5")
}

func pushToQueue(element string, clientConn net.Conn) {
	buffer.Lock()
	buffer.operations = append(buffer.operations, element)
	buffer.conns = append(buffer.conns, clientConn)
	buffer.Unlock()
}
