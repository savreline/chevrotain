package main

import (
	"context"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"../../util"

	"github.com/savreline/GoVector/govec"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	posCollection = "kvsp"
	negCollection = "kvsn"
)

// Global variables
var no int
var noStr string
var port string
var delay int   // TODO
var eLog string // TODO (1/2)
var noReplicas int
var curTick = 1 // TODO
var conns []*rpc.Client
var logger *govec.GoLog
var db *mongo.Database
var verbose = true     // TODO (1/2)
var flag = []int{0, 0} // TODO
var ticks [][]int      // TODO

// RPCExt is the RPC object that receives commands from the driver
type RPCExt int

// RPCInt is the RPC Object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, starts up the RPC server
func main() {
	var err error

	/* Parse args, initialize data structures */
	noReplicas, err = strconv.Atoi(os.Args[1])
	no, err = strconv.Atoi(os.Args[2])
	noStr = os.Args[2]
	port = os.Args[3]
	dbPort := os.Args[4]
	delay, err = strconv.Atoi(os.Args[5])
	conns = make([]*rpc.Client, noReplicas)
	ticks = make([][]int, noReplicas)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	/* Connect to MongoDB */
	dbClient, _ := util.Connect(noStr, dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB on "+dbPort)

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := util.CvRecord{Name: "Keys", Timestamp: logger.GetCurrentVC(), Values: []util.ValueEntry{}}
	_, err = db.Collection(posCollection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	_, err = db.Collection(negCollection).InsertOne(context.TODO(), newRecord)
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

	/* Start Server */
	util.PrintMsg(noStr, "RPC Server Listening on "+port)
	go rpc.Accept(l)
	select {}
}

// ConnectReplica connects this replica to others
func (t *RPCExt) ConnectReplica(args *util.RPCExtArgs, reply *int) error {
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

// TerminateReplica saves the logs to disk
func (t *RPCExt) TerminateReplica(args *util.RPCExtArgs, reply *int) error {
	if verbose == true {
		err := ioutil.WriteFile("Repl"+noStr+".txt", []byte(eLog), 0644)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}
	return nil
}
