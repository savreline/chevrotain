package cmrdt

import (
	"context"
	"fmt"

	"../util"
	"github.com/DistributedClocks/GoVector/govec"
	"go.mongodb.org/mongo-driver/bson"
)

/**********************/
/*** 1A: INSERT KEY ***/
/**********************/

// InsertKey inserts the given key with an empty array for values
func (t *RPCCmd) InsertKey(args *KeyArgs, reply *int) error {
	InsertKeyLocal(args.Key)
	InsertKeyGlobal(args.Key)
	return nil
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
	err := clients[0].Call("RPCObj.InsertKeyRPC", KeyArgs{key}, &result)
	if err != nil {
		util.PrintErr(err)
	}
	fmt.Println("Result from RPC", result)
}
