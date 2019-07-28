package main

import "net"

var TID = "000"

// var ipServer = "172.22.158.45"
// var portSlice = []string{"5251", "5252", "5253", "5254", "5255"}
var serverAddrs = []string{
	"172.22.158.45:5251",
	"172.22.94.54:5251",
	"172.22.156.46:5251",
	"172.22.158.46:5251",
	"172.22.94.55:5251",
}
var serverConnMap = make(map[string]net.Conn)
var coordinatorConn net.Conn

// var coordinatorAddress = "************"
//coordinator 的地址
var coordinatorAddress = "172.22.158.45:6666"
var beginFlag bool
