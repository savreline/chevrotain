package main

import (
	"strconv"
	"time"

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
		go simpleTest(i)
	}
	select {}
}

// simpleTest
func simpleTest(no int) {
	/* Connect to the Replica and Connect the Replica */
	var result int
	conn := util.RPCClient("DRIVER", ports[no])
	err := conn.Call("RPCExt.ConnectReplica", util.RPCExtArgs{}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Inserts */
	k := 0
	for i := 0; i < 20; i++ {
		key := (no+1)*100 + i
		conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 10; j++ {
			val := (no+1)*1000 + k
			conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < 20; i++ {
		key := (no+1)*100 + i
		for j := 0; j < 10; j++ {
			val := (no+1)*1000 + k
			conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Terminate */
	time.Sleep(3 * time.Second)
	err = conn.Call("RPCExt.TerminateReplica", util.RPCExtArgs{}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	util.PrintMsg("DRIVER", "Done on "+ports[no])
}
