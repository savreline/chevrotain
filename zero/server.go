package main

import (
	"net"
	"net/rpc"
	"os"
	"strconv"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"
)

// Constants
const (
	sCollection = "kvs"
)

// Global variables
var no int
var noStr string
var ports []string
var ips []string
var eLog string
var verbose bool        // print to info console?
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var delay int // emulated link delay

// RPCExt is the RPC object that receives commands from the client
type RPCExt int

// RPCInt is the RPC object for internal replica-to-replica communication
type RPCInt int

// Dummy variables to call RPCInt methods "locally" to insert data into the local database
var result int
var rpcint RPCInt

// Makes connection to the database, starts up the RPC server
func main() {
	var err, err1, err4 error

	/* Parse command line arguments */
	no, err1 = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port := os.Args[2]
	dbPort := os.Args[3]
	delay, err4 = strconv.Atoi(os.Args[4])
	if os.Args[5] == "v" {
		verbose = true
	} else {
		verbose = false
	}
	if err1 != nil {
		util.PrintErr(noStr, "CmdLine: no conversion", err)
	}
	if err4 != nil {
		util.PrintErr(noStr, "CmdLine: delay conversion", err)
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

	/* Start server */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)
	select {}
}

// InitReplica sets connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, ips[i], port)
	}
	return nil
}

// TerminateReplica in this implementation is just a place holder
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	return nil
}
