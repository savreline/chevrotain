package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	"../../util"
)

func simpleTest(no int, noKeys int, noVals int) {
	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no])
	var result int
	k := 0
	t := time.Now().UnixNano()

	/* Inserts */
	for i := 0; i < noKeys; i++ {
		key := (no+1)*100 + i
		latencies[no][k] = time.Now().UnixNano()
		calls[no][k] = conn.Go("RPCExt.InsertKey",
			util.RPCExtArgs{Key: strconv.Itoa(key)},
			&result, nil)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		k++

		for j := 0; j < noVals; j++ {
			val := (no+1)*1000 + k
			latencies[no][k] = time.Now().UnixNano()
			calls[no][k] = conn.Go("RPCExt.InsertValue",
				util.RPCExtArgs{Key: strconv.Itoa(key),
					Value: strconv.Itoa(val)},
				&result, nil)
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++
		}
	}

	delta1 := time.Now().UnixNano() - t

	for i, call := range calls[no] {
		if call != nil {
			<-call.Done
			latencies[no][i] = time.Now().UnixNano() - latencies[no][i]
		}
	}

	delta2 := time.Now().UnixNano() - t

	/* Terminate */
	util.Terminate(ports[no], conn, 5)

	/* Compute Average */
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

	/* Write Latencies to CSV */
	err := ioutil.WriteFile("Latencies"+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Print Latency */
	avg := float32(sum) / 1000000 / float32(k)
	avgT := float32(delta2) / 1000000 / float32(k)
	util.PrintMsg("DRIVER:", "Time Elapsed to send ops is (us): "+fmt.Sprint(float32(delta1)/float32(1000)))
	util.PrintMsg("DRIVER:", "Average Time per op to "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avgT))
	util.PrintMsg("DRIVER:", "Average latency to "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avg))
}
