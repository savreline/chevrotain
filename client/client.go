package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

// Constants
const (
	sCollection = "kvs"
	noKeys      = 210
	noVals      = 5
)

// Global variables
var ports []string
var ips []string
var noPerRepl int
var delay int         // delays between sending commands
var timeInt int       // time interval to initialize the replica with
var removes = false   // true if client is to test removes
var mongotest = false // true if client is to test mongoDb's replication
var term = false      // true if the replica is to be terminated after the test
var cnt int           // operation counter
var db *mongo.Database

// Map of latencies, associated wait group and lock
var latencies map[int]int64
var lock sync.Mutex
var wg sync.WaitGroup

// Wait group of the main method
var wgMain sync.WaitGroup

func main() {
	/* Parse command line arguments */
	var err error
	delay, err = strconv.Atoi(os.Args[1])
	timeInt, err = strconv.Atoi(os.Args[2])
	if os.Args[3] == "y" {
		mongotest = true
		dbClient, _ := util.ConnectDb("1", "localhost", "27018")
		db = dbClient.Database("chev")
		util.CreateCollection(db, "CLIENT", sCollection)
		util.PrintMsg("CLIENT", "Connected to DB")
	}
	if os.Args[4] == "y" {
		removes = true
	}
	if os.Args[5] == "y" {
		term = true
	}
	if err != nil {
		util.PrintErr("CLIENT", "CmdLine", err)
	}

	/* Init data structures */
	latencies = make(map[int]int64)

	/* Parse group membership */
	ips, ports, _, err = util.ParseGroupMembersCVS("../ports.csv", "")
	if err != nil {
		util.PrintErr("CLIENT", "GroupInfo", err)
	}
	noReplicas := len(ports)
	noPerRepl = noKeys / noReplicas

	/* Tests */
	for i := 0; i < noReplicas; i++ {
		wgMain.Add(1)
		go maintest(i)
	}
	for i := 0; i < noReplicas; i++ {
		// wgMain.Add(1)
		// go quicktest(i)
	}
	// startingPages := []string{"Java", "C++", "C--"}
	for i := 0; i < noReplicas; i++ {
		// wgMain.Add(1)
		// go wikiTest(startingPages[i], i)
	}
	wgMain.Wait()

	/* Process collected performance data */
	calcPerf()
	fmt.Println("Count: ", cnt)
}

// send the command
func sendCmd(key string, val string, cmdType util.OpCode, conn *rpc.Client) {
	wg.Add(1)
	defer wg.Done()
	var result int

	/* Send command and record time */
	t := time.Now().UnixNano()
	if cmdType == util.IK && !mongotest {
		conn.Call("RPCExt.InsertKey", util.RPCExtArgs{Key: key}, &result)
	} else if cmdType == util.IV && !mongotest {
		conn.Call("RPCExt.InsertValue", util.RPCExtArgs{Key: key, Value: val}, &result)
	} else if cmdType == util.RK && !mongotest {
		conn.Call("RPCExt.RemoveKey", util.RPCExtArgs{Key: key}, &result)
	} else if cmdType == util.RV && !mongotest {
		conn.Call("RPCExt.RemoveValue", util.RPCExtArgs{Key: key, Value: val}, &result)
	} else if cmdType == util.IK {
		util.InsertSKey(db.Collection(sCollection), "CLIENT", key)
	} else if cmdType == util.IV {
		util.InsertSValue(db.Collection(sCollection), "CLIENT", key, val, false)
	} else if cmdType == util.RK {
		util.RemoveSKey(db.Collection(sCollection), "CLIENT", key)
	} else {
		util.RemoveSValue(db.Collection(sCollection), "CLIENT", key, val)
	}

	/* Record latency and print progress to console */
	lock.Lock()
	cnt++
	latencies[cnt] = time.Now().UnixNano() - t
	if cnt%100 == 0 {
		fmt.Println("Current: ", cnt)
	}
	lock.Unlock()
}

// process collected performance data
func calcPerf() {
	/* Compute average */
	var sum int64
	for key := range latencies {
		sum += latencies[key]
	}

	/* Print latency */
	avg := float32(sum) / 1000000 / float32(cnt)
	util.PrintMsg("CLIENT", "Average latency is (ms): "+fmt.Sprint(avg))
}
