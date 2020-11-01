package main

import (
	"../cmrdt"
)

func main() {
	cmrdt.InitReplica(true, "27018", "8000", "1")
	cmrdt.ConnectReplica()
	cmrdt.InsertKey("1")
	// cmrdt.TerminateReplica()
	for {
	}
}
