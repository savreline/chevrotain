package main

import (
	"net/rpc"
	"os"
	"strconv"

	"../../util"
)

// Global variables
var calls [][]*rpc.Call
var latencies []map[int]int64
var ports []string
var delay int

func main() {
	/* Parse Command Line Arguments */
	var err error
	delay, err = strconv.Atoi(os.Args[1])
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Parse Group Membership */
	ports, _, err = util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	noReplicas := len(ports)

	/* Make Latency Maps */
	calls = make([][]*rpc.Call, noReplicas)
	latencies = make([]map[int]int64, noReplicas)
	for i := 0; i < noReplicas; i++ {
		latencies[i] = make(map[int]int64)
		calls[i] = make([]*rpc.Call, 1050)
	}

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		go simpleTest(i, 5, 2)
	}
	for i := 0; i < noReplicas; i++ {
		// go removeTest(i)
	}
	// wikiTest()

	select {}
}
