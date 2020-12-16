package util

import (
	"context"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitArgs are the arguments to Init/Terminate RPCExt Calls
type InitArgs struct {
	Bias    [2]bool
	TimeInt int
}

// RPCExtArgs are the arguments to any other RPCExt Call
type RPCExtArgs struct {
	Key, Value string
}

// ConnectDb to MongoDB on the given port, as per https://www.mongodb.com/golang
func ConnectDb(no string, ip string, port string) (*mongo.Client, context.Context) {
	urlString := "mongodb://" + ip + ":" + port + "/"

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		PrintErr(no, "MongoConn[FindClient]", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		PrintErr(no, "MongoConn[Connect]", err)
	}

	return client, ctx
}

// ConnectLocalDb connects a replica to the local Db for smoke test purposes
func ConnectLocalDb() *mongo.Database {
	noStr := "1"
	dbPort := "27018"

	/* Connect to MongoDB */
	dbClient, _ := ConnectDb(noStr, "localhost", dbPort)
	db := dbClient.Database("chev")
	PrintMsg(noStr, "Connected to DB on "+dbPort)
	return db
}

// ConnectClient connects driver to a replica
func ConnectClient(ip string, port string, t int) *rpc.Client {
	var result int
	conn := RPCClient("CLIENT", ip, port)
	err := conn.Call("RPCExt.InitReplica", InitArgs{Bias: [2]bool{true, true}, TimeInt: t}, &result)
	if err != nil {
		PrintErr("CLIENT", "ConnTo:"+port+" [InitReplica]", err)
	}
	return conn
}

// RPCClient makes an RPC connection to the given port
func RPCClient(no string, ip string, port string) *rpc.Client {
	dest := ip + ":" + port
	client, err := rpc.Dial("tcp", dest)
	if err != nil {
		PrintErr(no, "ConnTo:"+port+" [DialRPC]", err)
	}

	PrintMsg(no, "Connection made to "+dest)
	return client
}

// TerminateReplica is a command from the driver to terminate a replica
func TerminateReplica(port string, conn *rpc.Client, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
	var result int
	err := conn.Call("RPCExt.TerminateReplica", RPCExtArgs{}, &result)
	if err != nil {
		PrintErr("CLIENT", "ConnTo:"+port+" [Terminate]", err)
	}
	PrintMsg("CLIENT", "Done on "+port)
}
