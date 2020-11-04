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
	fmt.Printf("REPLICA "+strconv.Itoa(no+1)+": Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	// Insert entry for the given key
	newRecord := Record{key, []string{}}
	_, err = replicas[no].db.Collection("kvs").InsertOne(context.TODO(), newRecord)
	if err != nil {
		util.PrintErr(err)
	}
	util.PrintMsg(no, "Inserted Key "+key)
}

// InsertKeyGlobal broadcasts the insertKey operation to other replicas
func InsertKeyGlobal(key string, no int) {
	var result int
	var destNo int
	var flag = false

	for i, client := range replicas[no].clients {
		if i == no {
			flag = true
		}

		if client != nil {
			if flag {
				destNo = i + 1
			} else {
				destNo = i
			}
			fmt.Println("Sending RPC", no+1, "->", destNo+1)
			err := client.Call("RPCObj.InsertKeyRPC", KeyArgs{destNo, key}, &result)
			if err != nil {
				util.PrintErr(err)
			}
		}
	}

	util.PrintMsg(no, "Done Sending RPC Calls")
}
