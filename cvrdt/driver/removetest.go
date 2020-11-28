package main

import (
	"strconv"
	"time"

	"../../util"
)

func removeTest(no int) {
	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no])
	var result int

	/* Inserts: insert 2 keys */
	k := 0
	for i := 0; i < 2; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 0; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Removes: remove 1 key */
	k = 0
	for i := 0; i < 1; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.RemoveKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < 0; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.RemoveValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Merge */
	time.Sleep(5 * time.Second)
	err := conn.Call("RPCExt.GetCurrentSnapShot", util.RPCExtArgs{}, &result)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Terminate */
	util.Terminate(ports[no], conn)
}
