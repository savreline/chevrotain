package cmrdt

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"../util"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vrpc"
	"go.mongodb.org/mongo-driver/mongo"
)

var no string
var srvPort string
var db *mongo.Database
var ctx context.Context
var dbClient *mongo.Client
var clients []*rpc.Client
var logs []*govec.GoLog
var srvLogger *govec.GoLog

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

/********************************/
/*** 0: INIT, CONNECT REPLICA ***/
/********************************/

// InitReplica makes connection to the database, starts up the RPC server
func InitReplica(flag bool, iNo string, noReplicas int64, port string, dbPort string) {
	no = iNo
	clients = make([]*rpc.Client, noReplicas)
	logs = make([]*govec.GoLog, noReplicas)

	/* Connect to MongoDB */
	srvPort = port
	dbClient, ctx = util.Connect(dbPort)
	db = dbClient.Database("chev")
	util.PrintMsg(no, "Connected to DB")

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}

	/* Start Server */
	srvChanel := make(chan bool)
	go rpcserver(srvChanel)
	<-srvChanel
}

// RPCCmd is the RPC object that receives commands from the test application
type RPCCmd int

// ConnectReplica connects this replica to others
func (t *RPCCmd) ConnectReplica(args *ConnectArgs, reply *int) error {
	/* Parse Group Members */
	ports, err := util.ParseGroupMembersText("ports.txt", srvPort)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("REPLICA "+no+": Being Connected to, ports ", ports)

	/* Make RPC Connections */
	for i, port := range ports {
		rpcChan := make(chan *rpc.Client)
		logChan := make(chan *govec.GoLog)
		go util.RPCClient(rpcChan, logChan, port, "REPLICA "+no+": ")
		clients[i] = <-rpcChan
		logs[i] = <-logChan
	}

	return nil
}

// TerminateReplica closes the db connection
func TerminateReplica() {
	dbClient.Disconnect(ctx)
}

/*********************/
/*** 1: RPC SERVER ***/
/*********************/

func rpcserver(srvChanel chan bool) {
	/* Init RPC */
	util.PrintMsg(no, "Staring Server")
	server := rpc.NewServer()
	rpcobj := new(RPCObj)
	rpccmd := new(RPCCmd)
	server.Register(rpcobj)
	server.Register(rpccmd)
	l, e := net.Listen("tcp", ":"+srvPort)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	util.PrintMsg(no, "Listening on "+srvPort)

	/* Init cector clocks */
	srvLogger = govec.InitGoVector("server"+no, "server"+no, govec.GetDefaultConfig())
	options := govec.GetDefaultLogOptions()

	/* Acknowledge Readiness */
	util.PrintMsg(no, "RPC Server Ready")
	srvChanel <- true
	vrpc.ServeRPCConn(server, l, srvLogger, options)
}
