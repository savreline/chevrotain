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
var client *rpc.Client
var clLogger *govec.GoLog
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
func InitReplica(flag bool, no string, port string, dbPort string) {
	/* Connect to MongoDB */
	srvPort = port
	dbClient, ctx = util.Connect(dbPort)
	db = dbClient.Database("chev")
	fmt.Println("STATUS: Connected to DB")

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

// ConnectReplica connects this replica to others
func ConnectReplica() {
	/* Parse Group Members */
	ports, err := util.ParseGroupMembersText("ports.txt", srvPort)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println(ports)

	/* Make RPC Connections */
	for _, port := range ports {
		clChanel := make(chan bool)
		go rpcclient(clChanel, port)
		<-clChanel
	}
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
	fmt.Println("STATUS: Staring Server")
	rpcobj := new(RPCObj)
	server := rpc.NewServer()
	server.Register(rpcobj)
	l, e := net.Listen("tcp", ":"+srvPort)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("STATUS: Listening on " + srvPort)

	/* Init and Use Vector Clocks */
	srvLogger = govec.InitGoVector("server"+no, "server"+no, govec.GetDefaultConfig())
	options := govec.GetDefaultLogOptions()
	fmt.Println("STATUS: RPC Ready")
	srvChanel <- true
	vrpc.ServeRPCConn(server, l, srvLogger, options)
}

/*********************/
/*** 2: RPC CLIENT ***/
/*********************/

func rpcclient(clChanel chan bool, port string) {
	fmt.Println("STATUS: Staring Client: " + no)
	clLogger = govec.InitGoVector("client"+no, "client"+no, govec.GetDefaultConfig())
	options := govec.GetDefaultLogOptions()
	fmt.Println("STATUS: Client Clocks: " + no)

	var err error
	client, err = vrpc.RPCDial("tcp", "127.0.0.1:"+port, clLogger, options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("STATUS: Client Started: " + no)
	clChanel <- true
}
