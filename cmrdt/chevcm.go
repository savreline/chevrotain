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
	"go.mongodb.org/mongo-driver/bson"
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
func InitReplica(flag bool, dbPort string, port string, no string) {
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
	go rpcserver()
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
	go rpcclient()
}

// TerminateReplica closes the db connection
func TerminateReplica() {
	dbClient.Disconnect(ctx)
}

/**********************/
/*** 1A: INSERT KEY ***/
/**********************/

// InsertKey inserts the given key with an empty array for values
func InsertKey(key string) {
	InsertKeyLocal(key)
	InsertKeyGlobal(key)
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	srvLogger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	updateResult, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	newRecord := Record{key, []string{}}
	_, err = db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Inserted key", key)
}

// InsertKeyGlobal broadcasts the insertKey operation to other replicas
func InsertKeyGlobal(key string) {
	var result int
	err := client.Call("RPCObj.InsertKeyRPC", KeyArgs{key}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Result from RPC", result)
}

/************************/
/*** 2A: INSERT VALUE ***/
/************************/

// InsertValue inserts value into the given key
func InsertValue(key string, value string) {
	InsertValueLocal(key, value)
	InsertValueGlobal(key, value)
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string) {
	srvLogger.LogLocalEvent("Inserting value"+value, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: value}}}}

	updateResult, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)
}

// InsertValueGlobal broadcasts the insertValue operation to other replicas
func InsertValueGlobal(key string, value string) {
	var result int
	err := client.Call("InsertValueRPC", ValueArgs{key, value}, &result)
	if err != nil {
		util.PrintErr(err)
	}
}

/**********************/
/*** 3: RPC METHODS ***/
/**********************/

// RPCObj is the RPC Object
type RPCObj int

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	key string
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	key   string
	value string
}

// InsertKeyRPC receives incoming insert key call
func (t *RPCObj) InsertKeyRPC(args *KeyArgs, reply *int) error {
	fmt.Println("RPC Insert Key")
	*reply = 100
	return nil
}

// InsertValueRPC receives incoming insert value call
func (t *RPCObj) InsertValueRPC(args *ValueArgs, reply *int) error {
	fmt.Println("RPC Insert Value")
	return nil
}

/*********************/
/*** 4: RPC SERVER ***/
/*********************/

func rpcserver() {
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
	srvLogger = govec.InitGoVector("Server"+no, "LogFile"+no, govec.GetDefaultConfig())
	options := govec.GetDefaultLogOptions()
	fmt.Println("STATUS: RPC Ready")
	vrpc.ServeRPCConn(server, l, srvLogger, options)
}

/*********************/
/*** 5: RPC CLIENT ***/
/*********************/

func rpcclient() {
	fmt.Println("STATUS: Staring Client")
	clLogger = govec.InitGoVector("client", "clientlogfile", govec.GetDefaultConfig())
	options := govec.GetDefaultLogOptions()
	fmt.Println("STATUS: Client Clocks")
	var err error
	client, err = vrpc.RPCDial("tcp", "127.0.0.1:"+srvPort, clLogger, options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("STATUS: Client Started")
}
