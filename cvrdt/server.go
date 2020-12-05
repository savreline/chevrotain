package main

import (
	"context"
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
	posCollection  = "kvsp"
	negCollection  = "kvsn"
	permCollection = "kvs"
)

// Global variables
var no int
var noStr string
var port string
var ports []string
var noReplicas int
var eLog string
var verbose = true      // Print to info console?
var clock = 0           // Lamport clock: tick on broadcast and every local db op
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database

// Settings: bias towards add or removes for keys and values
// Settings: time interval between state exchanges
var bias [2]bool
var timeInt int

// Channels that activate the state exchange and garbage collection
// processes when the replica is initialized
var chanSE = make(chan bool)
var chanGC = make(chan bool)

// Current safe clock tick agreed upon by all replicas
var curSafeTick = 0

// Emulated link delay
var delay int

// RPCExt is the RPC object that receives commands from the driver
type RPCExt int

// RPCInt is the RPC Object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, starts up the RPC server
func main() {
	var err error

	/* Parse command link arguments */
	no, err = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port = os.Args[2]
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
	noReplicas = len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)

	/* Connect to MongoDB */
	dbClient, _ := util.ConnectDb(noStr, dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)

	/* Pre-allocate keys document */
	doc := util.CvDoc{Key: "Keys", Values: []util.CvRecord{}}
	_, err = db.Collection(posCollection).InsertOne(context.TODO(), doc)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	_, err = db.Collection(negCollection).InsertOne(context.TODO(), doc)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Init RPC */
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	rpc.Register(rpcint)
	rpc.Register(rpcext)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Start background processes */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)
	go runSE()
	go runGC()
	select {}
}

// InitReplica connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	bias = args.Bias
	timeInt = args.TimeInt

	/* Activate background processes */
	chanSE <- true
	chanGC <- true

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
