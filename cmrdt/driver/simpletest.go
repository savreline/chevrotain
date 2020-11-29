package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	"../../util"
)

func simpleTest(no int, noKeys int, noVals int, noop bool) {
	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no])
	var result int
	k := 0

	/* Send ops only to replica 0 */
	if !noop || (noop && no == 0) {
		/* Inserts */
		for i := 0; i < noKeys; i++ {
			key := (no+1)*100 + i
			t := time.Now().UnixNano()
			err := conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
			if err != nil {
				util.PrintErr("DRIVER", err)
			}
			latencies[no][k] = time.Now().UnixNano() - t
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++
			// }

			// for i := 0; i < noKeys; i++ {
			// 	key := (no+1)*100 + i
			for j := 0; j < noVals; j++ {
				val := (no+1)*1000 + k
				t := time.Now().UnixNano()
				err := conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
				if err != nil {
					util.PrintErr("DRIVER", err)
				}
				latencies[no][k] = time.Now().UnixNano() - t
				time.Sleep(time.Duration(delay) * time.Millisecond)
				k++
			}
		}
	}

	/* Terminate */
	if noop {
		util.Terminate(ports[no], conn, 30)
	} else {
		util.Terminate(ports[no], conn, 10)
	}

	/* Write Latencies to CSV */
	var str string
	var sum int64
	keys := make([]int, 0, len(latencies[no]))
	for key := range latencies[no] {
		keys = append(keys, key)
		sum += latencies[no][key]
	}
	sort.Ints(keys)
	for _, key := range keys {
		str = str + strconv.Itoa(key) + "," + strconv.FormatInt(latencies[no][key], 10) + "\n"
	}
	err := ioutil.WriteFile("Latencies"+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}
	avg := float32(sum) / 1000000 / float32(k)
	util.PrintMsg("DRIVER:", "Average Latency to Replica "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avg))
}
