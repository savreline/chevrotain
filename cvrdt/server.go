package main

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"../util"

	"go.mongodb.org/mongo-driver/mongo"
)

// Constants
const (
	posCollection = "kvsp"
	negCollection = "kvsn"
	sCollection   = "kvs"
)

// Global variables
var no int
var noStr string
var ports []string
var eLog string
var verbose = true      // print to info console?
var clock = 0           // lamport clock: tick on broadcast and every local db op
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var delay int // emulated link delay

// Settings: bias towards add or removes for keys and values
// Settings: time interval between state updates
var bias [2]bool
var timeInt int

// Channel that activate the state updates processes once
// the replica has been initialized
var chanSU = make(chan bool)

// Current safe clock tick agreed upon by all replicas
var curSafeTick = 0

// RPCExt is the RPC object that receives commands from the client
type RPCExt int

// RPCInt is the RPC object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, initializes data structures, starts up the RPC server
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
	ports, _, err = util.ParseGroupMembersCVS("../ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	noReplicas := len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)

	/* Connect to MongoDB */
	dbClient, _ := util.ConnectDb(noStr, dbPort)
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
	go runSU()
	select {}
}

// InitReplica sets the given parameters, activates background processes
// and connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	bias = args.Bias
	timeInt = args.TimeInt

	/* Activate background process */
	chanSU <- true

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, port)
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
