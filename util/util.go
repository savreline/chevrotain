package util

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/rpc"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SRecord is a static Db record
type SRecord struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// DDoc is a dynamic Db document
type DDoc struct {
	Key    string    `json:"key"`
	Values []DRecord `json:"values"`
}

// DRecord is a dynamic Db record (id is timestamp in the case of CvRDT)
type DRecord struct {
	Value string `json:"value"`
	ID    int    `json:"id"`
}

// InitArgs are the arguments to Init/Terminate RPCExt Calls
type InitArgs struct {
	Bias    [2]bool
	TimeInt int
}

// RPCExtArgs are the arguments to any other RPCExt Call
type RPCExtArgs struct {
	Key, Value string
}

// OpCode is an operation code
type OpCode int

// OpCodes
const (
	IK OpCode = iota + 1
	IV
	RK
	RV
	NO
)

// LookupOpCode translate operation code from string to op code
func LookupOpCode(opCode OpCode, noStr string) string {
	if opCode == IK {
		return "IK"
	} else if opCode == IV {
		return "IV"
	} else if opCode == RK {
		return "RK"
	} else if opCode == RV {
		return "RV"
	} else if opCode == NO {
		return "No-Op"
	} else {
		PrintErr(noStr, errors.New("LookupOpCode: unknown operation"))
		return ""
	}
}

// ConnectDb to MongoDB on the given port, as per https://www.mongodb.com/golang
func ConnectDb(no string, port string) (*mongo.Client, context.Context) {
	urlString := "mongodb://localhost:" + port + "/"

	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		PrintErr(no, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		PrintErr(no, err)
	}

	return client, ctx
}

// ConnectLocalDb connects a replica to the local Db for testing purposes
func ConnectLocalDb() *mongo.Database {
	noStr := "1"
	dbPort := "27018"

	/* Connect to MongoDB */
	dbClient, _ := ConnectDb(noStr, dbPort)
	db := dbClient.Database("chev")
	PrintMsg(noStr, "Connected to DB on "+dbPort)
	return db
}

// ParseGroupMembersCVS parses the supplied CVS group member file
func ParseGroupMembersCVS(file string, port string) ([]string, []string, error) {
	// adapted from https://stackoverflow.com/questions/24999079/reading-csv-file-in-go
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	clPorts := []string{}
	dbPorts := []string{}

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return clPorts, dbPorts, nil
		}

		/* Remove own port from results if appropriate */
		if row[0] != port {
			clPorts = append(clPorts, row[0])
			dbPorts = append(dbPorts, row[1])
		}
	}
}

// RPCClient makes an RPC connection to the given port
func RPCClient(no string, port string) *rpc.Client {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		PrintErr(no, err)
	}

	PrintMsg(no, "Connection made to "+port)
	return client
}

// ConnectClient connects driver to a replica
func ConnectClient(port string, t int) *rpc.Client {
	var result int
	conn := RPCClient("CLIENT", port)
	err := conn.Call("RPCExt.InitReplica", InitArgs{Bias: [2]bool{true, true}, TimeInt: t}, &result)
	if err != nil {
		PrintErr("CLIENT", err)
	}
	return conn
}

// TerminateReplica is a command from the driver to terminate a replica
func TerminateReplica(port string, conn *rpc.Client, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
	var result int
	err := conn.Call("RPCExt.TerminateReplica", RPCExtArgs{}, &result)
	if err != nil {
		PrintErr("CLIENT", err)
	}
	PrintMsg("CLIENT", "Done on "+port)
}

// DownloadDState downloads the contents of any dynamic collection
func DownloadDState(col *mongo.Collection, who string, drop string) []DDoc {
	var res []DDoc

	/* Download all key docs */
	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}

	/* Save downloaded info into res */
	if err = cursor.All(context.TODO(), &res); err != nil {
		PrintErr(who, err)
	}

	/* Drop the collection if asked */
	if drop == "1" {
		col.Drop(context.TODO())
	}
	return res
}

// DownloadSState downloads the contents of any static collection
func DownloadSState(col *mongo.Collection, who string, drop string) []SRecord {
	var result []SRecord

	opts := options.Find().SetSort(bson.D{{Key: "key", Value: 1}})
	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		PrintErr(who, err)
	}
	if err = cursor.All(context.TODO(), &result); err != nil {
		PrintErr(who, err)
	}
	if drop == "1" {
		col.Drop(context.TODO())
	}
	return result
}

// PrintDState prints a dynamic state to the console
func PrintDState(state []DDoc) {
	for _, doc := range state {
		fmt.Println(doc)
	}
	fmt.Println()
}

// PrintSState prints a static state to the console
func PrintSState(state []SRecord) {
	for _, record := range state {
		fmt.Println(record)
	}
	fmt.Println()
}

// PrintMsg prints message to console from a replica
func PrintMsg(no string, msg string) {
	if no == "CLIENT" || no == "TESTER" {
		fmt.Println(no + ": " + msg)
	} else {
		fmt.Println("REPLICA " + no + ": " + msg)
	}
}

// PrintErr prints error to console from a replica and exits
func PrintErr(no string, err error) {
	if no == "CLIENT" || no == "TESTER" {
		fmt.Println(no+": ", err)
	} else {
		fmt.Println("REPLICA "+no+": ", err)
	}
	os.Exit(1)
}

// GetRand generates a random number to emulate connection delays
func GetRand(no int) int {
	// https://golang.cafe/blog/golang-random-number-generator.html
	rand.Seed(time.Now().UnixNano())
	min := int(0.8 * float32(no))
	max := int(1.2 * float32(no))
	res := rand.Intn(max-min+1) + min
	return res
}

// EmulateDelay emulates link delay in all RPC responses
func EmulateDelay(delay int) {
	if delay > 0 {
		time.Sleep(time.Duration(GetRand(delay)) * time.Millisecond)
	}
}

// Max returns the maximum of a and b
func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
