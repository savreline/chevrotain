package main

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/savreline/GoVector/govec"
)

// Constants
const (
	sCollection = "kvs"
)

// Global variables
var no int
var noStr string
var noReplicas int
var ports []string
var ips []string
var eLog string
var verbose int         // print info console and to logs?
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var logger *govec.GoLog
var delay int // emulated link delay

// Settings: bias towards add or removes for keys and values
// Settings: minimum time interval between calls to process the queue
var bias [2]bool
var minTimeInt int

// Flag that activates the background process that periodically
// sends no-ops to other replicass once the replica has been initialized,
// along with a flag which indicates if a no-op has been sent
var flagNO = false
var sent = false

// Head of the operation queue, the associated lock and flag
// that activates periodic processing of the queue
var queue *ListNode
var lock sync.Mutex
var flagPQ = false
var queueLen = 0
var maxQueueLen int
var timeInt int
var times []int // slice of last 10 time intervals between queue processings
var lastT int64 // last time the queue was processed

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
	var err, err1, err4, err5, err6 error

	/* Parse command line arguments */
	no, err1 = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port := os.Args[2]
	dbPort := os.Args[3]
	delay, err4 = strconv.Atoi(os.Args[4])
	verbose, err5 = strconv.Atoi(os.Args[5])
	maxQueueLen, err6 = strconv.Atoi(os.Args[6])
	if err1 != nil {
		util.PrintErr(noStr, "CmdLine: no conversion", err)
	}
	if err4 != nil {
		util.PrintErr(noStr, "CmdLine: delay conversion", err)
	}
	if err5 != nil {
		util.PrintErr(noStr, "CmdLine: verbose conversion", err)
	}
	if err6 != nil {
		util.PrintErr(noStr, "CmdLine: maxQueueLen conversion", err)
	}

	/* Parse group member information */
	ips, ports, _, err = util.ParseGroupMembersCVS("../ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, "GroupInfo", err)
	}
	noReplicas = len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)
	ticks = make([][]int, noReplicas)

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Connect to MongoDB, Init collections (for performance) */
	dbClient, _ := util.ConnectDb(noStr, "localhost", dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)
	util.CreateCollection(db, noStr, sCollection)

	/* Init RPC */
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	rpc.Register(rpcint)
	rpc.Register(rpcext)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		util.PrintErr(noStr, "InitRPC", err)
	}

	/* Start server and background processes */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)
	go runNO()
	go runPQ()

	/* Save logs in case of Ctrl+C */
	go func() { // https://stackoverflow.com/questions/8403862/do-actions-on-end-of-execution
		channel := make(chan os.Signal)
		signal.Notify(channel, os.Interrupt)
		<-channel
		var result int
		rpcext.TerminateReplica(&util.RPCExtArgs{}, &result)
		os.Exit(0)
	}()
	select {}
}

// InitReplica sets the given parameters, activates background no-op sending processes
// and connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	bias = args.Bias
	minTimeInt = args.TimeInt
	timeInt = minTimeInt
	lastT = time.Now().UnixNano()

	/* Activate background processes */
	flagNO = true
	flagPQ = true

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, ips[i], port)
	}
	return nil
}

// TerminateReplica saves the logs to disk and stops background processes
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if verbose > 0 {
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, "WriteELog", err)
		}
	}
	eLog = ""

	/* Stop background processes */
	flagNO = false
	flagPQ = false
	return nil
}
