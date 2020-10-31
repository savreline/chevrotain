package main

import (
	"time"

	"../cmrdt"
)

func main() {
	cmrdt.InitReplica(true, "27018", "8000", "1")
	cmrdt.ConnectReplica()
	time.Sleep(5 * time.Second)
	cmrdt.InsertKey("1")
	// cmrdt.TerminateReplica()
	for {
	}
}
