package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"sort"
	"strconv"
	"time"

	"../../util"
)

func simpleTest(no int, noKeys int, noVals int, noop bool) {
	/* Connect to the Replica and Connect the Replica */
	conn := util.ConnectDriver(ports[no])
	k := 0
	t := time.Now().UnixNano()

	/* Send ops only to replica 0 */
	if !noop || (noop && no == 0) {

		/* Inserts */
		for i := 0; i < noKeys; i++ {
			key := (no+1)*100 + i
			wg.Add(1)
			go insertKey(key, k, no, conn)
			time.Sleep(time.Duration(delay) * time.Millisecond)
			k++

			for j := 0; j < noVals; j++ {
				val := (no+1)*1000 + k
				wg.Add(1)
				go insertValue(key, val, k, no, conn)
				time.Sleep(time.Duration(delay) * time.Millisecond)
				k++

				if k%100 == 0 {
					fmt.Println("K: ", k)
				}
			}
		}
	}

	delta1 := time.Now().UnixNano() - t
	fmt.Println("Waiting")
	wg.Wait()

	/* Terminate */
	if noop {
		util.Terminate(ports[no], conn, 30)
	} else {
		util.Terminate(ports[no], conn, 5)
	}

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
		num := int(float32(latencies[no][key]) / 1000000)
		str = str + strconv.Itoa(key) + "," + strconv.Itoa(num) + "\n"
	}

	/* Write Latencies to CSV */
	err := ioutil.WriteFile("Latencies"+strconv.Itoa(no)+".csv", []byte(str), 0644)
	if err != nil {
		util.PrintErr("DRIVER", err)
	}

	/* Print Latency */
	avg := float32(sum) / 1000000 / float32(k)
	// avgT := float32(delta2) / 1000000 / float32(k)
	util.PrintMsg("DRIVER:", "Time Elapsed to send ops is (us): "+fmt.Sprint(float32(delta1)/float32(1000)))
	// util.PrintMsg("DRIVER:", "Average Time per op to "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avgT))
	util.PrintMsg("DRIVER:", "Average latency to "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avg))
}

func insertKey(key int, k int, no int, conn *rpc.Client) {
	defer wg.Done()
	var result int
	t := time.Now().UnixNano()
	conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: strconv.Itoa(key)}, &result)
	lock.Lock()
	latencies[no][k] = time.Now().UnixNano() - t
	lock.Unlock()
}

func insertValue(key int, val int, k int, no int, conn *rpc.Client) {
	defer wg.Done()
	var result int
	t := time.Now().UnixNano()
	conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: strconv.Itoa(key), Value: strconv.Itoa(val)}, &result)
	lock.Lock()
	latencies[no][k] = time.Now().UnixNano() - t
	lock.Unlock()
}
