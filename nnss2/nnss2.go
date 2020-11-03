package main

import (
	"fmt"

	"../cmrdt"
	"../util"
)

// REPLICA 2: 8002, 27019
func main() {
	cmrdt.InitReplica(true, "2", "8002", "27019")

	fmt.Println("Replica Initialized, proceed to connecting? [y]")
	var ans string
	_, err := fmt.Scanln(&ans)
	if err != nil {
		util.PrintErr(err)
	}

	cmrdt.ConnectReplica()
	cmrdt.InsertKey("2")
	// cmrdt.TerminateReplica()
	for {
	}
}
