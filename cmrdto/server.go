package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"../util"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/savreline/GoVector/govec"
)

// Constants
const (
	dCollection = "kvsd"
	sCollection = "kvs"
)

// Global variables
var no int
var noStr string
var ports []string
var ips []string
var eLog string
var iLog string
var verbose int         // print to info console?
var conns []*rpc.Client // RPC connections to other replicas
var db *mongo.Database
var logger *govec.GoLog
var id int    // unique ids associated with elements
var delay int // emulated link delay

// Slice of channels that are used for communication with waiting RPC Calls
// and the associated lock
var chans []chan bool
var lock sync.Mutex

// RPCExt is the RPC object that receives commands from the client
type RPCExt int

// RPCInt is the RPC object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, initializes data structures, starts up the RPC server
func main() {
	var err, err1, err4, err5 error

	/* Parse command line arguments */
	no, err1 = strconv.Atoi(os.Args[1])
	noStr = os.Args[1]
	port := os.Args[2]
	dbPort := os.Args[3]
	delay, err4 = strconv.Atoi(os.Args[4])
	verbose, err5 = strconv.Atoi(os.Args[5])
	if err1 != nil {
		util.PrintErr(noStr, "CmdLine: no conversion", err)
	}
	if err4 != nil {
		util.PrintErr(noStr, "CmdLine: delay conversion", err)
	}
	if err5 != nil {
		util.PrintErr(noStr, "CmdLine: verbose conversion", err)
	}

	/* Parse group member information */
	ips, ports, _, err = util.ParseGroupMembersCVS("../ports.csv", port)
	if err != nil {
		util.PrintErr(noStr, "GroupInfo", err)
	}
	noReplicas := len(ports) + 1

	/* Init data structures */
	conns = make([]*rpc.Client, noReplicas)
	id = no * 100000

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Connect to MongoDB, Init collections (for performance) */
	dbClient, _ := util.ConnectDb(noStr, "localhost", dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)
	util.CreateCollection(db, noStr, dCollection)
	util.CreateCollection(db, noStr, sCollection)

	/* Init RPC */
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	rpc.Register(rpcint)
	rpc.Register(rpcext)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		util.PrintErr(noStr, "RPCInit", err)
	}

	/* Start server */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)

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

// InitReplica connects this replica to others
func (t *RPCExt) InitReplica(args *util.InitArgs, reply *int) error {
	for i, port := range ports {
		conns[i] = util.RPCClient(noStr, ips[i], port)
	}
	return nil
}

// TerminateReplica generates the "lookup" view collection of the database
// and saves the logs to disk
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if verbose > 0 {
		eLog = eLog + "\nFinal Clock:\n" + fmt.Sprint(logger.GetCurrentVC())
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, "WriteELog", err)
		}
	}
	if verbose > 0 {
		err := ioutil.WriteFile("iRepl"+noStr+".txt", []byte(iLog), 0644)
		if err != nil {
			util.PrintErr(noStr, "WriteILog", err)
		}
	}
	return nil
}
