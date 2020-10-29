package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"../util"

	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database
var logger *govec.GoLog
var clPort string
var srvPort string
var clConn []*net.UDPConn
var srvConn net.PacketConn

// Record is a DB Record
type Record struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func main() {
	/* Connect to MongoDB, reading port number from the command line */
	dbPort := os.Args[1]
	clPort = os.Args[2]
	srvPort = os.Args[3]
	client, ctx := util.Connect(dbPort)
	db = client.Database("chev")
	defer client.Disconnect(ctx)

	/* Initialize GoVector logger */
	logger = govec.InitGoVector("MyProcess", "LogFile", govec.GetDefaultConfig())

	/* Pre-allocate Keys entry */
	newRecord := Record{"Keys", []string{}}
	_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Pre-allocated keys")

	/* Parse Group Members */
	ports, err := util.ParseGroupMembers("ports.csv", clPort, srvPort)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println(ports)

	/* Setup Connection */
	setupConnection(ports)

	// /* Tests */
	// InsertKey("1")
	// InsertKey("2")
	// InsertValue("1", "Hello")
	// InsertValue("2", "Bye")
}

/****************************/
/*** 1: INSERT/REMOVE KEY ***/
/****************************/

// InsertKey inserts the given key with an empty array for values
func InsertKey(key string) {
	InsertKeyLocal(key)
	InsertKeyGlobal(key)
}

// RemoveKey removes the given key
func RemoveKey(key string) {
	// TODO
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string) {
	logger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
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

}

/******************************/
/*** 2: INSERT/REMOVE VALUE ***/
/******************************/

// InsertValue inserts value into the given key
func InsertValue(key string, value string) {
	InsertValueLocal(key, value)
	InsertValueGlobal(key, value)
}

// RemoveValue removes value from the given key
func RemoveValue(key string, value string) {
	// TODO
}

// InsertValueLocal inserts the value into the local db
func InsertValueLocal(key string, value string) {
	logger.LogLocalEvent("Inserting value"+value, govec.GetDefaultLogOptions())
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

}

/*************************/
/*** 3: UDP CONNECTION ***/
/*************************/

// from https://github.com/DistributedClocks/GoVector/blob/master/example/ClientServer/ClientServer.go
func setupConnection(ports map[string]string) {
	// set up server
	// fmt.Println("Listening on server....")
	// srvConn, err := net.ListenPacket("udp", ":"+srvPort)
	// if err != nil {
	// 	util.PrintErr(err)
	// }

	for aClPort, aSrvPort := range ports {
		// resolve addresses
		rAddr, errR := net.ResolveUDPAddr("udp4", ":"+aClPort)
		fmt.Println(rAddr)
		util.PrintErr(errR)
		lAddr, errL := net.ResolveUDPAddr("udp4", ":"+aSrvPort)
		fmt.Println(lAddr)
		util.PrintErr(errL)

		// make connection
		// aClConn, errDial := net.DialUDP("udp", lAddr, rAddr)
		// clConn = append(clConn, aClConn)
		// util.PrintErr(errDial)
	}
}
