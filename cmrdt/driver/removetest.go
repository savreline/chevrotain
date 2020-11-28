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

	/* Inserts: insert 2 keys, 1 with values */
	k := 0
	for i := 1; i < 2; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	for i := 0; i < 1; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)

		for j := 0; j < 2; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++
		}
	}

	/* Removes: remove 1 key and 1 value from other key */
	k = 0
	for i := 1; i < 2; i++ {
		key := (no+1)*100 + i
		err := conn.Call("RPCExt.RemoveKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
		if err != nil {
			util.PrintErr("DRIVER", err)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	for i := 0; i < 1; i++ {
		key := (no+1)*100 + i
		for j := 0; j < 1; j++ {
			val := (no+1)*1000 + k
			err := conn.Call("RPCExt.RemoveValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++
		}
	}

	/* Terminate */
	util.Terminate(ports[no], conn, 3)
}
