package main

import (
	"fmt"

	"../cmrdt"
	"../util"
)

// REPLICA 1: 8001, 27018
func main() {
	cmrdt.InitReplica(true, "1", "8001", "27018")

	fmt.Println("Replica Initialized, proceed to connecting? [y]")
	var ans string
	_, err := fmt.Scanln(&ans)
	if err != nil {
		util.PrintErr(err)
	}

	cmrdt.ConnectReplica()
	cmrdt.InsertKey("1")
	// cmrdt.TerminateReplica()
	for {
	}
}
