package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectArgs are the arguments to the ConnectReplica call (a dummy structure)
type ConnectArgs struct {
}

// KeyArgs are the arguments to the InsertKeyRPC call
type KeyArgs struct {
	Key string
}

// ValueArgs are the arguments to the InsertValueRPC call
type ValueArgs struct {
	Key, Value string
}

// Connect to MongoDB on the given port, as per https://www.mongodb.com/golang
func Connect(port string) (*mongo.Client, context.Context) {
	urlString := "mongodb://localhost:" + port + "/"

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

// ParseGroupMembersCVS parses the supplied CVS group member file
func ParseGroupMembersCVS(file string, port string) ([]string, []string, error) {
	// from https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	ports := []string{}
	dbPorts := []string{}

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return ports, dbPorts, nil
		}

		if row[0] != port {
			ports = append(ports, row[0])
			dbPorts = append(dbPorts, row[1])
		}
	}
}

// RPCClient makes an RPC connection
func RPCClient(port string, who string) *rpc.Client {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(who + "Connection made to " + port)
	return client
}

// PrintMsg prints message to console from a replica
func PrintMsg(no int, msg string) {
	fmt.Println("REPLICA " + strconv.Itoa(no) + ": " + msg)
}

// PrintErr prints error
// from https://github.com/DistributedClocks/GoVector/blob/master/example/ClientServer/ClientServer.go
func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
