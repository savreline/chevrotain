package main

import (
	"context"
	"fmt"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// inserts the record in either positive collection (add) or negative collection (remove)
func insertLocalRecord(key string, value string, collection string, record *util.DRecord) {
	/* In no ready to go record is supplied, tick the clock and make one,
	otherwise check if an exact identical entry already exists */
	if record == nil {
		record = &util.DRecord{Value: value, ID: clock}
		lastRPC = time.Now().UnixNano()
		printTime = true
	} else {
		var res util.DDoc
		filter := bson.D{{Key: "key", Value: key},
			{Key: "values", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "value", Value: record.Value},
					{Key: "id", Value: record.ID}}}}}}
		err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
		if err == nil { // found the record, no need to do anything
			return
		}
	}

	/* Tick the clock */
	clock++

	/* Check if the document to be updated exists, if not, make one */
	var dbResult util.DRecord
	filter := bson.D{{Key: "key", Value: key}}
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&dbResult)
	if err != nil {
		keyEntry := &util.DDoc{Key: key, Values: []util.DRecord{}}
		_, err := db.Collection(collection).InsertOne(context.TODO(), keyEntry)
		if err != nil {
			util.PrintErr(noStr, "I-L:"+key+":"+value+" [Find]", err)
		}
	}

	/* Update the document */
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "values", Value: record}}}}
	_, err = db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, "I-L:"+key+":"+value+" [Update]", err)
	}

	/* Print to console */
	if verbose > 1 {
		if collection == posCollection && key == "Keys" {
			util.PrintMsg(noStr, "Inserted Key "+value)
		} else if collection == negCollection && key == "Keys" {
			util.PrintMsg(noStr, "Removed Key "+value)
		} else if collection == posCollection {
			util.PrintMsg(noStr, "Inserted Value "+value+" on key "+key)
		} else {
			util.PrintMsg(noStr, "Removed Value "+value+" on key "+key)
		}
	}

	/* Log this operation */
	if verbose > 0 {
		iLog = iLog + collection + ":" + key + ":" + value + ":" + fmt.Sprint(record.ID) + "\t"
	}

	count++
	if !gc && count == TOTALOPS {
		// https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
		curSafeTick = int(^uint(0) >> 1)
		fCount = count
		mergeCollections()
		util.PrintMsg(noStr, "Received all operations after (s): "+
			fmt.Sprint(float32(time.Now().UnixNano()-lastRPC)/float32(1000000000))+":"+fmt.Sprint(count))
	}
}

// insert key into the static collection
func insertKey(key string) {
	util.InsertSKey(db.Collection(sCollection), noStr, key)
}

// insert value into the static collection
func insertValue(key string, value string) {
	util.InsertSValue(db.Collection(sCollection), noStr, key, value)
}

// removes key from the static collection
func removeKey(key string) {
	util.RemoveSKey(db.Collection(sCollection), noStr, key)
}

// removes value from the static collection
func removeValue(key string, value string) {
	util.RemoveSValue(db.Collection(sCollection), noStr, key, value)
}

// prints a dynamic state to the log
func printDState(state []util.DDoc, name string) string {
	res := name + "\n"
	for _, doc := range state {
		res = res + fmt.Sprint(doc) + "\n"
	}
	res = res + "\n"
	return res
}

// prints a static state to the log
func printSState(state []util.SRecord) string {
	res := "STATIC\n"
	for _, record := range state {
		res = res + fmt.Sprint(record) + "\n"
	}
	res = res + "\n"
	return res
}

// log an insertion or a removal
func printToLog(msg string) {
	if verbose > 0 {
		iLog = iLog + msg + "\n"
	}
}
