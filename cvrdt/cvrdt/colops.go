package main

import (
	"context"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mergeCollection(state []util.CvRecord, collection string) {

	for _, record := range state {

		/* Keys: filter is set to key's name, if not found, insert record */
		var res util.CvRecord
		filter := bson.D{{Key: "name", Value: record.Name}}
		err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
		if err != nil {
			InsertKeyLocal(record.Name, collection, &record)
		}

		/* Values: look for with the key, if not found, insert vrecord */
		for _, vrecord := range record.Values {
			var res util.ValueEntry
			filter := bson.D{{Key: "name", Value: record.Name},
				{Key: "values", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{{Key: "value", Value: vrecord.Value}}}}}}
			err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
			if err != nil {
				InsertValueLocal(record.Name, vrecord.Value, collection, &vrecord)
			}
		}
	}
}

// merge positive and negative collections
func mergeCollections() {
	/* Iterate over the positive collection */
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := db.Collection(posCollection).Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		util.PrintErr(noStr, err)
	}

	for cursor.Next(context.TODO()) {
		/* Since item been in posCollection, default action is to insert it */
		var insert = true

		/* Get next item */
		var posRecord util.CvRecord
		if err := cursor.Decode(&posRecord); err != nil {
			util.PrintErr(noStr, err)
		}

		/* Try to locate the item in the negative collection */
		var negRecord util.CvRecord
		var filter = bson.D{{Key: "name", Value: posRecord.Name}}
		err := db.Collection(negCollection).FindOne(context.TODO(), filter).Decode(&negRecord)

		if err == nil {
			/* If found, compare time stamps and decide */
			cmp := posRecord.Timestamp.CompareClocks(negRecord.Timestamp)

			/* If negative record is the last one, do not insert the item */
			if cmp == 2 {
				insert = false
			}

			/* If concurrent, decide as per rules */
			if cmp == 1 {
				if settings[0] == 0 {
					insert = false
				}
			}

			/* Delete negative record */
			_, err = db.Collection(negCollection).DeleteOne(context.TODO(), filter)
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}

		/* Insert if need be */
		if insert == true {
			newRecord := util.CmRecord{Name: posRecord.Name, Values: []string{}}
			_, err := db.Collection("kvs").InsertOne(context.TODO(), newRecord)
			if err != nil {
				util.PrintErr(noStr, err)
			}
		}

		/* Delete positive record */
		_, err = db.Collection(posCollection).DeleteOne(context.TODO(), filter)
		if err != nil {
			util.PrintErr(noStr, err)
		}
	}
}
