package main

import (
	"net"
	"sync"
)

//coordinator 的地址
var coordinatorAddress = "172.22.158.45:6666"
var coordinatorConn net.Conn
var serverPort = "5251"
var committedMap = make(map[string]string)             // {object : value}
var tentativeList = make(map[string]map[string]string) // {tid : {object : value}}
var commitStatus = make(map[string]bool)               // {tid : commit?}
var tidMap = make(map[string]net.Conn)                 // {tid : net.Conn}

var serverLock = struct {
	sync.RWMutex
	//key:=object, value:=locktype
	lockInfo map[string]int
	//key:=object, value:=holder or holders
	holder map[string][]string
}{
	lockInfo: make(map[string]int),
	holder:   make(map[string][]string),
}

var buffer = struct {
	sync.RWMutex
	operations []string
	conns      []net.Conn
}{
	operations: make([]string, 0),
	conns:      make([]net.Conn, 0),
}
