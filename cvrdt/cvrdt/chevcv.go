package cvrdt

/* In this file:
0. Definitions of Replica, Record, ValueEntry structs
	Definitions of ConnectArgs, RPCExt
1. Init, InitReplica (connets to Db, initializes keys entry, starts up RPC server)
2. ConnectReplica (makes connections to other replicas)
3. TerminateReplica (disconnect from Db)
*/

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"

	"../../util"

	"github.com/savreline/GoVector/govec"
	"github.com/savreline/GoVector/govec/vclock"
	"github.com/savreline/GoVector/govec/vrpc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	posCollection = "kvsp"
	negCollection = "kvsn"
)

// Replica holds parameters associated with a replica
type Replica struct {
	port     string
	db       *mongo.Database
	ctx      context.Context
	dbClient *mongo.Client
	clients  []*rpc.Client
	logger   *govec.GoLog
}

var noReplicas int
var replicas []Replica

// Record is a DB Record
type Record struct {
	Name   string        `json:"name"`
	Time   vclock.VClock `json:"time"`
	Values []ValueEntry  `json:"values"`
}

// ValueEntry is a value along with the timestamp
type ValueEntry struct {
	Value string        `json:"name"`
	Time  vclock.VClock `json:"time"`
}

// ConnectArgs are the arguments to the ConnectReplica call
type ConnectArgs struct {
	No int
}

// Init initializes chevcv
func Init(iNoReplicas int, flag bool, time int) {
	replicas = make([]Replica, iNoReplicas)
	noReplicas = iNoReplicas
}

// InitReplica makes connection to the database, starts up the RPC server
func InitReplica(no int, port string, dbPort string) {
	clients := make([]*rpc.Client, noReplicas)
	noStr := strconv.Itoa(no + 1)

	/* Connect to MongoDB */
	dbClient, ctx := util.Connect(dbPort)
	db := dbClient.Database("chev")
	util.PrintMsg(no, "Connected to DB")

	/* Init vector clocks */
	logger := govec.InitGoVector("R"+noStr, "R"+noStr, govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", logger.GetCurrentVC(), []ValueEntry{}}
	_, err := db.Collection(posCollection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	_, err = db.Collection(negCollection).InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}

	/* Start Server */
	channel := make(chan bool)
	go rpcserver(channel, logger, no, port)
	<-channel

	replicas[no] = Replica{port, db, ctx, dbClient, clients, logger}
}

// RPCExt is the RPC object that receives commands from the driver
type RPCExt int

// ConnectReplica connects this replica to others
func (t *RPCExt) ConnectReplica(args *ConnectArgs, reply *int) error {
	no := args.No
	noStr := strconv.Itoa(no + 1)

	/* Parse Group Members */
	ports, _, err := util.ParseGroupMembersCVS("ports.csv", replicas[no].port)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("REPLICA "+noStr+": Being Connected to ports", ports)

	/* Make RPC Connections */
	for i, port := range ports {
		channel := make(chan *rpc.Client)
		go util.RPCClient(channel, replicas[no].logger, port, "REPLICA "+noStr+": ")
		replicas[no].clients[i] = <-channel
	}

	return nil
}

// TerminateReplica closes the db connection
func (t *RPCExt) TerminateReplica(args *ConnectArgs, reply *int) error {
	no := args.No
	replicas[no].dbClient.Disconnect(replicas[no].ctx)
	return nil
}

// RPC Server
func rpcserver(srvChanel chan bool, logger *govec.GoLog, no int, port string) {
	/* Init RPC */
	util.PrintMsg(no, "Staring Server")
	server := rpc.NewServer()
	rpcint := new(RPCInt)
	rpcext := new(RPCExt)
	server.Register(rpcint)
	server.Register(rpcext)
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	util.PrintMsg(no, "Listening on "+port)

	/* Acknowledge Readiness */
	options := govec.GetDefaultLogOptions()
	util.PrintMsg(no, "RPC Server Ready")
	srvChanel <- true

	vrpc.ServeRPCConn(server, l, logger, options)
}
