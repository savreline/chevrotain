package main

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"

	"../../util"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
)

// Constants
const (
	collectionName = "kvs"
)

// Global variables
var no int
var noStr string
var port string
var eLog string
var iLog string
var noReplicas int
var conns []*rpc.Client
var logger *govec.GoLog
var db *mongo.Database
var chans = make(map[chan vclock.VClock]chan vclock.VClock)
var channel = make(chan int)
var lock sync.Mutex
var verbose = true

// RPCExt is the RPC object that receives commands from the driver
type RPCExt int

// RPCInt is the RPC Object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, starts up the RPC server
func main() {
	var err error

	/* Parse args, initialize data structures */
	noReplicas, err = strconv.Atoi(os.Args[1])
	no, _ = strconv.Atoi(os.Args[2])
	noStr = os.Args[2]
	port = os.Args[3]
	dbPort := os.Args[4]
	conns = make([]*rpc.Client, noReplicas)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Connect to MongoDB */
	dbClient, _ := util.ConnectDb(noStr, dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Init RPC */
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	rpc.Register(rpcext)
	rpc.Register(rpcint)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Start Server */
	util.PrintMsg(noStr, "Server Listening on "+port)
	go rpc.Accept(l)
	select {}
}

// ConnectReplica connects this replica to others
func (t *RPCExt) ConnectReplica(args *util.InitArgs, reply *int) error {
	/* Set up args */
	// channel <- args.TimeInt

	/* Parse Group Members */
	ports, _, err := util.ParseGroupMembersCVS("../driver/ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, port)
	}

	return nil
}

// TerminateReplica writes to the log
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if verbose == true {
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, err)
		}
		err = ioutil.WriteFile("iRepl"+strconv.Itoa(no)+".txt", []byte(iLog), 0644)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func broadcastClockValue(clockValue vclock.VClock) {
	for _, channel := range chans {
		channel <- clockValue
	}
}
