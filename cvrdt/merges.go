package main

import (
	"context"
	"fmt"
	"time"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// merge either the positive or the negative collection during state updates
func mergeState(state []util.DDoc, collection string) {
	for _, doc := range state {
		for _, record := range doc.Values {
			insertLocalRecord(doc.Key, record.Value, collection, &record)
		}
	}
}

// merge positive and negative collections during garbage collection
func mergeCollections() {
	var posLen, negLen int
	if verbose > 0 {
		iLog = iLog + "\n\nMerging Collections on Tick: " + fmt.Sprint(curSafeTick) + "\n"
		iLog = iLog + printDState(util.DownloadDState(db, "REPLICA "+noStr, posCollection, "0"), "POSITIVE")
		iLog = iLog + printDState(util.DownloadDState(db, "REPLICA "+noStr, negCollection, "0"), "NEGATIVE")
		iLog = iLog + printSState(util.DownloadSState(db, "REPLICA "+noStr, "0"))
	}

	/* Download the positive collection and negative collections (for efficiency) */
	posState := util.DownloadDState(db, "REPLICA "+noStr, posCollection, "0")
	negState := util.DownloadDState(db, "REPLICA "+noStr, negCollection, "0")

	/* Iterate over documents in the positive collection */
	for _, posDoc := range posState { // util.DDoc
		if len(posDoc.Values) > 0 {
			posLen++
		}

		/* Look for the corresponding doc(s) in the negative collection */
		var negDoc, found = locateNegDoc(negState, posDoc.Key)

		/* If negative doc not found, just go ahead and insert all records */
		if !found {
			for _, record := range posDoc.Values {
				if record.ID > curSafeTick {
					continue
				}
				if posDoc.Key == "Keys" {
					printToLog("IK:" + record.Value + ":" + "NoNeg")
					insertKey(record.Value)
				} else {
					printToLog("IV:" + posDoc.Key + ":" + record.Value + ":" + "NoNeg")
					insertValue(posDoc.Key, record.Value)
				}
				deleteDRecord(posDoc.Key, record, posCollection)
			}
			continue
		}

		/* Iterate over records in the document */
		for i := 0; i < len(posDoc.Values); i++ {
			record := posDoc.Values[i]
			var insert = true
			if record.ID > curSafeTick {
				continue
			}

			/* Get max times of all identical elements in positive and negative collections;
			consider only elements below the current safe tick;
			remove elements as they been proceed */
			posTimestamp := getMaxTimestamp(posDoc.Values, record.Value)
			negTimestamp := getMaxTimestamp(negDoc.Values, record.Value)

			/* Determine if the element needs to be inserted */
			if posTimestamp < negTimestamp ||
				(posTimestamp == negTimestamp && posDoc.Key == "Keys" && bias[0]) ||
				(posTimestamp == negTimestamp && posDoc.Key != "Keys" && bias[1]) {
				insert = false
			}

			/* Delete all dynamic records as they been processed */
			deleteDRecord(posDoc.Key, record, posCollection)
			deleteDRecord(posDoc.Key, record, negCollection)

			/* Insert or delete into static collection as required */
			if insert && posDoc.Key == "Keys" {
				printToLog("IK:" + record.Value + ":" +
					fmt.Sprint(posTimestamp) + ":" + fmt.Sprint(negTimestamp))
				insertKey(record.Value)
			} else if insert && posDoc.Key != "Keys" {
				printToLog("IV:" + posDoc.Key + ":" + record.Value + ":" +
					fmt.Sprint(posTimestamp) + ":" + fmt.Sprint(negTimestamp))
				insertValue(posDoc.Key, record.Value)
			} else if !insert && posDoc.Key == "Keys" {
				printToLog("RK:" + record.Value + ":" +
					fmt.Sprint(posTimestamp) + ":" + fmt.Sprint(negTimestamp))
				removeKey(record.Value)
			} else {
				printToLog("RV:" + posDoc.Key + ":" + record.Value + ":" +
					fmt.Sprint(posTimestamp) + ":" + fmt.Sprint(negTimestamp))
				removeValue(posDoc.Key, record.Value)
			}
		}
	}

	/* Iterate over documents in the negative collection:
	those documents didn't have a corresponding positve entry
	and must be removed */
	negState = util.DownloadDState(db, "REPLICA "+noStr, negCollection, "0")
	for _, negDoc := range negState { // util.DDoc
		if len(negDoc.Values) > 0 {
			negLen++
		}

		for _, record := range negDoc.Values {
			if record.ID > curSafeTick {
				continue
			}
			if negDoc.Key == "Keys" {
				printToLog("RK:" + record.Value + ":" + "NoPos")
				removeKey(record.Value)
			} else {
				printToLog("RV:" + negDoc.Key + ":" + record.Value + ":" + "NoPos")
				removeValue(negDoc.Key, record.Value)
			}
			deleteDRecord(negDoc.Key, record, negCollection)
		}
	}

	if posLen == 0 && negLen == 0 && printTime {
		delta := float32(time.Now().UnixNano()-lastRPC) / float32(1000000000)
		util.PrintMsg(noStr, "mergeCollections reached zero dynamic state after (s): "+
			fmt.Sprint(delta))
		opts := options.Update().SetUpsert(true)
		db.Collection("time").UpdateOne(
			context.TODO(),
			bson.D{},
			bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: delta}}}},
			opts)
		printTime = false
	}
	if verbose > 0 {
		util.PrintMsg(noStr, "state lengths are "+fmt.Sprint(posLen)+":"+fmt.Sprint(negLen))
		iLog = iLog + "States After Merge on Tick: " + fmt.Sprint(curSafeTick) + "\n"
		iLog = iLog + printDState(util.DownloadDState(db, "REPLICA "+noStr, posCollection, "0"), "POSITIVE")
		iLog = iLog + printDState(util.DownloadDState(db, "REPLICA "+noStr, negCollection, "0"), "NEGATIVE")
		iLog = iLog + printSState(util.DownloadSState(db, "REPLICA "+noStr, "0"))
	}
}

// return the maximum timestamp of the given element
func getMaxTimestamp(arr []util.DRecord, val string) int {
	res := -1
	for _, record := range arr {
		if record.Value == val && record.ID <= curSafeTick && record.ID > res {
			res = record.ID
		}
	}
	return res
}

// locates all entries in the negative state corresponding to the given key
func locateNegDoc(negState []util.DDoc, key string) (util.DDoc, bool) {
	negDoc := util.DDoc{Key: key, Values: []util.DRecord{}}
	found := false
	for _, doc := range negState {
		if key == doc.Key {
			found = true
			for _, val := range doc.Values {
				negDoc.Values = append(negDoc.Values, val)
			}
		}
	}
	return negDoc, found
}

// deletes a processed recored from the dynamic database
func deleteDRecord(key string, record util.DRecord, collection string) {
	filter := bson.D{{Key: "key", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: bson.D{{
			Key: "value", Value: bson.D{{Key: "$eq", Value: record.Value}}}, {
			Key: "id", Value: bson.D{{Key: "$lte", Value: curSafeTick}}}}}}}}
	_, err := db.Collection(collection).UpdateMany(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, "Del-D:"+key+":"+record.Value, err)
	}
}
