package main

import (
	"os"
	"strconv"
	"sync"

	"../../util"
)

// Global variables
var latencies []map[int]int64
var ports []string
var delay int
var lock sync.Mutex
var wg sync.WaitGroup
var t int

func main() {
	/* Parse Group Membership */
	var err error
	delay, err = strconv.Atoi(os.Args[1])
	t, err = strconv.Atoi(os.Args[2])
	ports, _, err = util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	noReplicas := len(ports)

	/* Make Slice of Calls to Wait For */
	latencies = make([]map[int]int64, noReplicas)
	for i := 0; i < noReplicas; i++ {
		latencies[i] = make(map[int]int64)
	}

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		go simpleTest(i, 50, 20, false)
	}
	for i := 0; i < noReplicas; i++ {
		// go removeTest(i)
	}
	select {}
}
