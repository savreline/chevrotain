package main

import (
	"strconv"
	"time"

	"../../util"
)

func simpleTest(no int, keys int, vals int) {
	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no])
	var result int

	/* First Batch of Inserts */
	k := 0
	for i := 0; i < keys; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		for j := 0; j < vals; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Second Batch of Inserts: Values into pre-exisiting keys */
	time.Sleep(5000 * time.Millisecond)
	for i := 0; i < keys; i++ {
		key := (no+1)*100 + i
		for j := 0; j < vals; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			k++
		}
	}

	/* Terminate */
	util.Terminate(ports[no], conn, 3)
}
