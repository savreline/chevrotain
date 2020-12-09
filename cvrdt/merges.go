package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
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
	/* Download the positive collection and negative collections (for efficiency) */
	posState := util.DownloadDState(db.Collection(posCollection), "REPLICA "+noStr, "0")
	negState := util.DownloadDState(db.Collection(negCollection), "REPLICA "+noStr, "0")

	/* Iterate over documents in the positive collection */
	for _, posDoc := range posState { // util.DDoc

		/* Look for the corresponding doc in the negative collection */
		var negDoc util.DDoc
		var found = false
		for _, doc := range negState {
			if posDoc.Key == doc.Key {
				negDoc = doc
				found = true
			}
		}

		/* If negative doc not found, just go ahead and insert all records */
		if !found {
			for _, record := range posDoc.Values {
				if record.ID > curSafeTick {
					continue
				}
				if posDoc.Key == "Keys" {
					insertKey(record.Value)
				} else {
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
				insertKey(record.Value)
			} else if insert && posDoc.Key != "Keys" {
				insertValue(posDoc.Key, record.Value)
			} else if !insert && posDoc.Key == "Keys" {
				removeKey(record.Value)
			} else {
				removeValue(posDoc.Key, record.Value)
			}
		}

		/* Iterate over documents in the negative collection:
		those documents didn't have a corresponding positve entry
		and must be removed */
		for _, negDoc := range negState { // util.DDoc
			for _, record := range negDoc.Values {
				if record.ID > curSafeTick {
					continue
				}
				if negDoc.Key == "Keys" {
					removeKey(record.Value)
				} else {
					removeValue(record.Value, negDoc.Key)
				}
				deleteDRecord(negDoc.Key, record, negCollection)
			}
		}
	}
}

// return the maximum timestamp of the given element
func getMaxTimestamp(arr []util.DRecord, val string) int {
	res := -1
	for _, record := range arr {
		if record.Value == val && record.ID > res {
			res = record.ID
		}
	}
	return res
}

// deletes a processed recored from the dynamic database
func deleteDRecord(key string, record util.DRecord, collection string) {
	filter := bson.D{{Key: "key", Value: key}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "values", Value: bson.D{{
			Key: "value", Value: bson.D{{
				Key: "$eq", Value: record.Value}}}}}}}}
	_, err := db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.PrintErr(noStr, err)
	}
}
