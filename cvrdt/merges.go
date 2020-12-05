package main

import (
	"context"

	"../util"
	"go.mongodb.org/mongo-driver/bson"
)

// merge either the positive or the negative collection
func mergeState(state []util.CvDoc, collection string) {
	for _, doc := range state {
		for _, record := range doc.Values {
			InsertLocalRecord(doc.Key, record.Value, collection, &record)
		}
	}
}

// merge positive and negative collections
func mergeCollections() {
	var err error

	/* Download the positive collection and negative collections (for efficiency) */
	posState := util.DownloadCvState(db.Collection(posCollection), "REPLICA "+noStr, "0")
	negState := util.DownloadCvState(db.Collection(negCollection), "REPLICA "+noStr, "0")

	/* Iterate over documents in the positive collection */
	for posDoc := range posState { // util.CvDoc

		/* Look for the corresponding doc in the negative collection */
		var negDoc util.CvDoc
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
				if posDoc.Key == "Keys" {
					insertKey(record.Value)
				} else {
					insertValue(record.Value, posDoc.Key)
			}
			continue;
		}

		/* Iterate over records in the document */
		for i := 0; i < len(posDoc.Values); {
			record = posDoc.Values[i]
			var insert = true

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

			/* Delete all positive records as they been processed */
			filter := bson.D{{Key: "key", Value: posDoc.Key}}
			update := bson.D{{Key: "$pull", Value: bson.D{
				{Key: "values", Value: bson.D{{
					Key: "value", Value: bson.D{{
						Key: "$eq", Value: record.Value}}}}}}}}
			_, err := db.Collection(posCollection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				util.PrintErr(noStr, err)
			}

			/* Delete all positive records as they been processed */
			_, err = db.Collection(negCollection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				util.PrintErr(noStr, err)
			}

			/* Insert or delete as required */
			if insert && posDoc.Key == "Keys" {
				insertKey(record.Value)
			} else if insert && posDoc.Key != "Keys" {
				insertValue(record.Value, posDoc.Key)
			} else if !insert && posDoc.Key == "Keys" {
				removeKey(record.Value)
			} else {
				removeValue(record.Value, posDoc.Key)
			}
		}
	}

	/* Iterate over documents in the negative collection:
	those documents didn't have a corresponding positve entry
	and must be removed */
	for negDoc := range negState { // util.CvDoc
		for _, record := range negDoc.Values {
			if && negDoc.Key == "Keys" {
				removeKey(record.Value)
			} else {
				removeValue(record.Value, negDoc.Key)
			}
		}
	}
}

// return the maximum timestamp of the given element
func getMaxTimestamp(arr []util.CvRecord, val string) int {
	res := -1
	j := 0
	for i, record := range arr {
		if record.Value == val && record.Timestamp > res {
			res = record.Timestamp
			arr[j] = arr[i]
			j++
		}
	}
	return res
}
