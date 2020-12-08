package main

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/savreline/GoVector/govec"
)

// Constants
const (
	sCollection = "kvs"
	MAXQUEUELEN = 20
)

// Global variables
var no int
var noStr string
var noReplicas int
var ports []string
var ips []string
var eLog string
var verbose = true      // print to info console?
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var logger *govec.GoLog
var delay int // emulated link delay

// Settings: bias towards add or removes for keys and values
// Settings: time interval between state updates
var bias [2]bool
var timeInt int

// Channel that activates the background process that periodically
// sends no-ops to other replicass once the replica has been initialized,
// along with a flag which indicates if a no-op has been sent
var chanNO = make(chan bool)
var sent = false

// Head of the operation queue and the associated lock
var queue *ListNode
var lock sync.Mutex
var queueLen = 0

// Lists of clock ticks seen thus far from other replicas and
// the current safe clock tick
var ticks [][]int
var curSafeTick = 1

// RPCExt is the RPC object that receives commands from the client
type RPCExt int

// RPCInt is the RPC object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, starts up the RPC server
func main() {
	var err error

	/* Parse command link arguments */
	no, err = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port := os.Args[2]
	dbPort := os.Args[3]
	delay, err = strconv.Atoi(os.Args[4])
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Parse group member information */
	ips, ports, _, err = util.ParseGroupMembersCVS("../ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	noReplicas = len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)
	ticks = make([][]int, noReplicas)

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Connect to MongoDB */
	dbClient, _ := util.ConnectDb(noStr, "locahost", dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)

	/* Init RPC */
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	rpc.Register(rpcint)
	rpc.Register(rpcext)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Start server and background processes */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)
	go runNoOps()
	select {}
}

// InitReplica sets the given parameters, activates background no-op sending processes
// and connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	bias = args.Bias
	timeInt = args.TimeInt

	/* Activate background processes */
	chanNO <- true

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, ips[i], port)
	}
	return nil
}

// TerminateReplica saves the logs to disk
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if verbose {
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}
	return nil
}
