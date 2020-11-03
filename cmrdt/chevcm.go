package cmrdt

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"

	"../util"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vrpc"
	"go.mongodb.org/mongo-driver/mongo"
)

// Replica holds parameters associated with a replica
type Replica struct {
	port     string
	db       *mongo.Database
	ctx      context.Context
	dbClient *mongo.Client
	clients  []*rpc.Client
	loggers  []*govec.GoLog
	logger   *govec.GoLog
}

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

var noReplicas int
var replicas []Replica

/********************************/
/*** 0: INIT, CONNECT REPLICA ***/
/********************************/

// Init initializes chevcm
func Init(iNoReplicas int) {
	replicas = make([]Replica, iNoReplicas)
	noReplicas = iNoReplicas
}

// InitReplica makes connection to the database, starts up the RPC server
func InitReplica(flag bool, no int, port string, dbPort string) {
	clients := make([]*rpc.Client, noReplicas)
	loggers := make([]*govec.GoLog, noReplicas)
	noStr := strconv.Itoa(no)

	/* Connect to MongoDB */
	dbClient, ctx := util.Connect(dbPort)
	db := dbClient.Database("chev")
	util.PrintMsg(noStr, "Connected to DB")

	/* Init vector clocks */
	logger := govec.InitGoVector("server"+noStr, "server"+noStr, govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}

	/* Start Server */
	channel := make(chan bool)
	go rpcserver(channel, logger, noStr, port)
	<-channel

	replicas[no] = Replica{port, db, ctx, dbClient, clients, loggers, logger}
}

// RPCCmd is the RPC object that receives commands from the test application
type RPCCmd int

// ConnectReplica connects this replica to others
func (t *RPCCmd) ConnectReplica(args *ConnectArgs, reply *int) error {
	no := args.No
	noStr := strconv.Itoa(no)

	/* Parse Group Members */
	ports, err := util.ParseGroupMembersText("ports.txt", replicas[no].port)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("REPLICA "+noStr+": Being Connected to, ports", ports)

	/* Make RPC Connections */
	for i, port := range ports {
		rpcChan := make(chan *rpc.Client)
		logChan := make(chan *govec.GoLog)
		go util.RPCClient(rpcChan, logChan, port, "REPLICA "+noStr+": ")
		replicas[no].clients[i] = <-rpcChan
		replicas[no].loggers[i] = <-logChan
	}

	return nil
}

// TerminateReplica closes the db connection
func (t *RPCCmd) TerminateReplica(args *ConnectArgs, reply *int) error {
	no := args.No
	replicas[no].dbClient.Disconnect(replicas[no].ctx)
	return nil
}

/*********************/
/*** 1: RPC SERVER ***/
/*********************/

func rpcserver(srvChanel chan bool, logger *govec.GoLog, no string, port string) {
	/* Init RPC */
	util.PrintMsg(no, "Staring Server")
	server := rpc.NewServer()
	rpcobj := new(RPCObj)
	rpccmd := new(RPCCmd)
	server.Register(rpcobj)
	server.Register(rpccmd)
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
