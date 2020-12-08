package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"../util"
)

// Global variables
var ports []string
var delay int
var timeInt int    // time interval to initialize the replica with
var verbose = true // save detailed latency information to csv?

func main() {
	/* Parse command line arguments */
	var err error
	delay, err = strconv.Atoi(os.Args[1])
	timeInt, err = strconv.Atoi(os.Args[2])
	if err != nil {
		util.PrintErr("CLIENT", err)
	}

	/* Parse group membership */
	ports, _, err = util.ParseGroupMembersCVS("../ports.csv", "")
	if err != nil {
		util.PrintErr("CLIENT", err)
	}
	noReplicas := len(ports)

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		go test1(i, 5, 5)
	}
	// wikiTest()
	select {}
}

// send the command
func sendCmd(key string, val string, cnt int, cmdType util.OpCode,
	conn *rpc.Client, latencies map[int]int64, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	var result int

	/* Send command and record time */
	t := time.Now().UnixNano()
	if cmdType == util.IK {
		conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: key}, &result)
	} else if cmdType == util.IV {
		conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: key, Value: val}, &result)
	} else if cmdType == util.RK {
		conn.Call("RPCExt.RemoveKey", util.RPCExtArgs{Key: key}, &result)
	} else {
		conn.Call("RPCExt.RemoveValue", util.RPCExtArgs{Key: key, Value: val}, &result)
	}
	latencies[cnt] = time.Now().UnixNano() - t

	/* Print progress to console */
	if cnt%100 == 0 {
		fmt.Println("Current: ", cnt)
	}
}

// process collected performance data
func calcPerf(delta int64, cnt int, no int, latencies map[int]int64) {
	/* Compute average */
	var str string
	var sum int64
	keys := make([]int, 0, len(latencies))
	for key := range latencies {
		keys = append(keys, key)
		sum += latencies[key]
	}

	/* Print latency */
	avg := float32(sum) / 1000000 / float32(cnt)
	timeElps := float32(delta) / 1000000
	util.PrintMsg("CLIENT", "Time Elapsed to send ops is (ms): "+fmt.Sprint(timeElps))
	util.PrintMsg("CLIENT", "Average latency to "+strconv.Itoa(no)+" is (ms):"+fmt.Sprint(avg))

	/* Write latencies to CSV */
	if verbose {
		sort.Ints(keys)
		for _, key := range keys {
			num := int(float32(latencies[key]) / 1000000)
			str = str + strconv.Itoa(key) + "," + strconv.Itoa(num) + "\n"
		}
		err := ioutil.WriteFile("Latencies"+strconv.Itoa(no)+".csv", []byte(str), 0644)
		if err != nil {
			util.PrintErr("CLIENT", err)
		}
	}
}
