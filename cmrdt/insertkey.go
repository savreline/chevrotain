package cmrdt

import (
	"context"
	"fmt"
	"strconv"

	"../util"
	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

/**********************/
/*** 1A: INSERT KEY ***/
/**********************/

// InsertKey inserts the given key with an empty array for values
func (t *RPCCmd) InsertKey(args *KeyArgs, reply *int) error {
	InsertKeyLocal(args.Key, args.No)
	InsertKeyGlobal(args.Key, args.No)
	return nil
}

// InsertKeyLocal inserts the key into the local db
func InsertKeyLocal(key string, no int) {
	replicas[no].logger.LogLocalEvent("Inserting Key"+key, govec.GetDefaultLogOptions())
	filter := bson.D{{Key: "name", Value: "Keys"}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: key}}}}

	// Update global keys entry
	updateResult, err := replicas[no].db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	// Insert entry for the given key
	newRecord := Record{key, []string{}}
	_, err = replicas[no].db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	util.PrintMsg(strconv.Itoa(no), "Inserted Key "+key)
}

// InsertKeyGlobal broadcasts the insertKey operation to other replicas
func InsertKeyGlobal(key string, no int) {
	var result int
	err := replicas[no].clients[0].Call("RPCObj.InsertKeyRPC", KeyArgs{no, key}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Result from RPC", result)
}
