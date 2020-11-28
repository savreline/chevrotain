package main

import (
	"../../util"
)

// Global variables
var ports []string

func main() {
	/* Parse Group Membership */
	var err error
	ports, _, err = util.ParseGroupMembersCVS("ports.csv", "")
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	noReplicas := len(ports)

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		// go simpleTest(i, 2, 2)
	}
	for i := 0; i < noReplicas; i++ {
		go removeTest(i)
	}
	select {}
}
