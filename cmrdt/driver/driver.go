package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"../../util"
)

// Global variables
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
	latencies = make([]map[int]int64, noReplicas)
	for i := 0; i < noReplicas; i++ {
		latencies[i] = make(map[int]int64)
	}

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		go simpleTest(i, 50, 20)
	}
	for i := 0; i < noReplicas; i++ {
		// go removeTest(i)
	}
	// wikiTest()

	select {}
}

// https://golang.cafe/blog/golang-random-number-generator.html
func getRand() int {
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 100
	res := rand.Intn(max-min+1) + min
	return res
}
