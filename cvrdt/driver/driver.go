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
	// for i := 0; i < noReplicas; i++ {
	// 	go simpleTest(i)
	// }
	for i := 0; i < noReplicas; i++ {
		go deleteTest(i)
	}
	select {}
}

// simpleTest
func simpleTest(no int) {
	/* Connect to the Replica and Connect the Replica */
	var result int
	conn := util.RPCClient("DRIVER", ports[no])
	err := conn.Call("RPCExt.ConnectReplica", util.InitArgs{Settings: [2]int{0, 0}, TimeInt: 5000}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Inserts */
	k := 0
	for i := 0; i < 2; i++ {
		key := (no+1)*100 + i
		conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 1; j++ {
			val := (no+1)*1000 + k
			conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < 2; i++ {
		key := (no+1)*100 + i
		for j := 0; j < 1; j++ {
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

// deleteTest
func deleteTest(no int) {
	/* Connect to the Replica and Connect the Replica */
	var result int
	conn := util.RPCClient("DRIVER", ports[no])
	err := conn.Call("RPCExt.ConnectReplica", util.InitArgs{Settings: [2]int{0, 0}, TimeInt: 5000}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Inserts */
	k := 0
	for i := 0; i < 2; i++ {
		key := (no+1)*100 + i
		conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 0; j++ {
			val := (no+1)*1000 + k
			conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Removes */
	k = 0
	for i := 0; i < 1; i++ {
		key := (no+1)*100 + i
		conn.Call("RPCExt.RemoveKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 0; j++ {
			val := (no+1)*1000 + k
			conn.Call("RPCExt.RemoveValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Merge */
	time.Sleep(5 * time.Second)
	err = conn.Call("RPCExt.GetCurrentSnapShot", util.RPCExtArgs{}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Terminate */
	err = conn.Call("RPCExt.TerminateReplica", util.RPCExtArgs{}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	util.PrintMsg("DRIVER", "Done on "+ports[no])
}
