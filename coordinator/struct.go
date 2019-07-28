package main

import "net"

var coordinatorPort string

// var ipServer = "172.22.158.45"                                   //server的ip地址
// var portSlice = []string{"5251", "5252", "5253", "5254", "5255"} //server的监听端口
/*var serverAddrs = []string{
	"172.22.158.45:5251",
	"172.22.158.45:5252",
	"172.22.158.45:5253",
	"172.22.158.45:5254",
	"172.22.158.45:5255",
}
*/

var serverConnSet = make(map[net.Conn]bool) // {serverConn : no-use}
var tidSet = make(map[string]bool)          // {tid : no-use}
var canCommitCount = make(map[string]int)   // {tid : count}
var canAbortCount = make(map[string]bool)   // {tid : bool}
var tidNO = 1
var graphForever = Graph{}
