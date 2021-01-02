package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertKey inserts the key into the local db
func (t *RPCInt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	record := util.SRecord{Key: args.Key, Values: []string{}}
	_, err := db.Collection(sCollection).InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, "IK:"+args.Key, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Inserted Key "+args.Key)
	}
	return nil
}

// InsertValue inserts the given value into the local db
func (t *RPCInt) InsertValue(args *util.RPCExtArgs, reply *int) error {
	/* Define filters */
	filter := bson.D{{Key: "key", Value: args.Key}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: args.Value}}}}

	/* Do the update */
	_, err := db.Collection(sCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, "IV:"+args.Key+":"+args.Value, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Inserted Value "+args.Value)
	}
	return nil
}

// RemoveKey removes the given key from the local db
func (t *RPCInt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	filter := bson.D{{Key: "key", Value: args.Key}}
	_, err := db.Collection(sCollection).DeleteOne(context.TODO(), filter)
	if err != nil {
		util.PrintErr(noStr, "RK:"+args.Key, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Deleted Key "+args.Key)
	}
	return nil
}

// RemoveValue the given value from the local db
func (t *RPCInt) RemoveValue(args *util.RPCExtArgs, reply *int) error {
	/* Define filters */
	filter := bson.D{{Key: "key", Value: args.Key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: args.Value}}}}

	/* Do the delete */
	_, err := db.Collection(sCollection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, "RV:"+args.Key+":"+args.Value, err)
	}
	if verbose {
		util.PrintMsg(noStr, "Deleted Value "+args.Value)
	}
	return nil
}
