package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"

	"../util"

	"go.mongodb.org/mongo-driver/mongo"
)

// Constants
const (
	posCollection = "kvsp"
	negCollection = "kvsn"
	sCollection   = "kvs"
	TOTALOPS      = 1523
)

// Global variables
var no int
var noStr string
var ports []string
var ips []string
var eLog string
var iLog string
var verbose int         // print to info console?
var gc = false          // run with garbage collection?
var clock = 0           // Lamport clock: tick on broadcast and every local db op
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var delay int // emulated link delay

// Settings: bias towards add or removes for keys and values
// Settings: time interval between state updates
var bias [2]bool
var timeInt int

// Flag that activates the state update processes on the main replica
// once it has been initialized, along with the lock that must
// be acquired while merging states and/or merging collections
var flagSU = false
var lock sync.Mutex

// Current safe clock tick agreed upon by all replicas and
// the current safe clock tick on this replica
var curSafeTick = 0
var mySafeTick = 0

// time of the last incoming RPCExt call along with the flag which indicates
// whether to print time since last RPC call to console, and is reset to false
// once the lengths of positive and negative collections reach zero
var lastRPC int64
var printTime bool

// variables that keep count of total number of database operations for
// statistical purposes in the non-garbage collected version, fCount
// fixes the count once the last incoming RPC call has been received
var count int
var fCount int

// RPCExt is the RPC object that receives commands from the client
type RPCExt int

// RPCInt is the RPC object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, initializes data structures, starts up the RPC server
func main() {
	var err error

	/* Parse command line arguments */
	no, err = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port := os.Args[2]
	dbPort := os.Args[3]
	delay, err = strconv.Atoi(os.Args[4])
	verbose, err = strconv.Atoi(os.Args[5])
	if os.Args[6] == "y" {
		gc = true
	}
	if err != nil {
		util.PrintErr(noStr, "CmdLine", err)
	}

	/* Parse group member information */
	ips, ports, _, err = util.ParseGroupMembersCVS("../ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, "GroupInfo", err)
	}
	noReplicas := len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)

	/* Connect to MongoDB, Init collections (for performance) */
	dbClient, _ := util.ConnectDb(noStr, "localhost", dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)
	util.CreateCollection(db, noStr, posCollection)
	util.CreateCollection(db, noStr, negCollection)
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
	if no == 1 {
		go runSU()
	}
	select {}
}

// InitReplica sets the given parameters, activates background processes
// and connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	bias = args.Bias
	timeInt = args.TimeInt

	/* Activate background process */
	flagSU = true

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, ips[i], port)
	}
	return nil
}

// TerminateReplica saves the logs to disk
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if !gc {
		util.PrintMsg(noStr, "Counts are "+fmt.Sprint(count)+":"+fmt.Sprint(fCount))
	}
	if verbose > 0 {
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, "WriteELog", err)
		}
	}
	if verbose > 0 {
		err := ioutil.WriteFile("iRepl"+noStr+".txt", []byte(iLog), 0644)
		if err != nil {
			util.PrintErr(noStr, "WriteELog", err)
		}
	}
	return nil
}
