package main

/* In this file:
0. Globals: conns, logger, db, no, port
	Definitions: Record, RPCExt, RPCInt
1. Main (connet to db, init clocks, init keys entry, start up RPC server)
2. ConnectReplica (make conns to other replicas)
*/

import (
	"context"
	"io/ioutil"
	"log"
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

// Global variables
var no int
var noStr string
var port string
var eLog string
var conns []*rpc.Client
var logger *govec.GoLog
var db *mongo.Database
var chans = make(map[chan vclock.VClock]chan vclock.VClock)
var lock = &sync.Mutex{}
var verbose = false

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// RPCExt is the RPC object that receives commands from the driver
type RPCExt int

// RPCInt is the RPC Object for internal replica-to-replica communication
type RPCInt int

// Makes connection to the database, starts up the RPC server
func main() {
	/* Parse args, initialize data structures */
	noReplicas, _ := strconv.Atoi(os.Args[1])
	no, _ = strconv.Atoi(os.Args[2])
	noStr = os.Args[2]
	port = os.Args[3]
	dbPort := os.Args[4]
	conns = make([]*rpc.Client, noReplicas)

	/* Connect to MongoDB */
	dbClient, _ := util.Connect(dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(no, "Connected to DB on "+dbPort)

	/* Init vector clocks */
	logger = govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}

	/* Init RPC */
	server := rpc.NewServer()
	rpcext := new(RPCExt)
	rpcint := new(RPCInt)
	server.Register(rpcext)
	server.Register(rpcint)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	/* Start Server */
	util.PrintMsg(no, "Server Listening on "+port)
	go server.Accept(l)
	select {}
}

// ConnectReplica connects this replica to others
func (t *RPCExt) ConnectReplica(args *util.ConnectArgs, reply *int) error {
	/* Parse Group Members */
	ports, _, err := util.ParseGroupMembersCVS("../driver/ports.csv", port)
	if err != nil {
		util.PrintErr(err)
	}

	/* Make RPC Connections */
	for i, port := range ports {
		conns[i] = util.RPCClient(port, "REPLICA "+strconv.Itoa(no)+": ")
	}

	return nil
}

// TerminateReplica writes to the log
func (t *RPCExt) TerminateReplica(args *util.ConnectArgs, reply *int) error {
	if verbose == true {
		err := ioutil.WriteFile("Repl"+strconv.Itoa(no)+".txt", []byte(eLog), 0644)
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
