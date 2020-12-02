package main

import (
	"context"
	"strconv"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
)

func mergeCollection(state []util.CvRecord, collection string) {
	for _, entry := range state {
		for _, record := range entry.Values {
			var res util.ValueEntry
			filter := bson.D{{Key: "name", Value: entry.Name}, // could be Keys
				{Key: "values", Value: bson.D{ // look to see if we already have the value
					{Key: "$elemMatch", Value: bson.D{{Key: "value", Value: record.Value}}}}}}
			err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&res)
			if err != nil { // if we don't have it, need to insert the record
				InsertLocalRecord(entry.Name, record.Value, collection, &record)
			}
		}
	}
}

// merge positive and negative collections (need to refactor this to work symmetrically)
func mergeCollections() {
	/* Download the positive collection */
	cursor, err := db.Collection(posCollection).Find(context.TODO(), bson.D{})
	if err != nil {
		util.PrintErr(noStr, err)
	}

	for cursor.Next(context.TODO()) {
		/* Get the positive array */
		var posEntry, negEntry util.CvRecord
		if err := cursor.Decode(&posEntry); err != nil {
			util.PrintErr(noStr, err)
		}

		/* Try to find the negative array */
		filter := bson.D{{Key: "name", Value: posEntry.Name}}
		err = db.Collection(negCollection).FindOne(context.TODO(), filter).Decode(&negEntry)
		if err != nil {
			util.PrintErr(noStr, err)
		}

		/* Optimization: If not found, just copy the array directly */

		/* Record-by-record */
		for _, posRecord := range posEntry.Values {
			var insert = true

			/* Determine if this record is below the current safe tick */
			var skip = false
			for i := 0; i < noReplicas; i++ {
				if int(posRecord.Timestamp["R"+strconv.Itoa(i+1)]) > curLocalTick {
					skip = true
				}
			}

			if !skip {
				/* Try to find the record in the negative array */
				negRecord, found := locateValEntry(negEntry.Values, posRecord.Value)

				/* If found, decide on clocks */
				if found {
					cmp := posRecord.Timestamp.CompareClocks(negRecord.Timestamp)

					/* If negative record is the last one, do not insert the item */
					if cmp == 2 { // TODO: may have many records
						insert = false
					}

					/* If concurrent, decide as per rules */
					if cmp == 1 {
						if settings[0] == 0 {
							insert = false
						}
					}

					/* Delete positive record */
					update := bson.D{{Key: "$pull", Value: bson.D{
						{Key: "values", Value: bson.D{{
							Key: "value", Value: bson.D{{
								Key: "$eq", Value: posRecord.Value}}}}}}}}
					_, err := db.Collection(posCollection).UpdateOne(context.TODO(), filter, update)
					if err != nil {
						util.PrintErr(noStr, err)
					}

					/* Delete negative record */
					_, err = db.Collection(negCollection).UpdateOne(context.TODO(), filter, update)
					if err != nil {
						util.PrintErr(noStr, err)
					}
				}

				/* Insert new key into final collection (if need be) */
				if insert && posEntry.Name == "Keys" {
					var res util.CmRecord
					filter := bson.D{{Key: "name", Value: posRecord.Value}}
					err := db.Collection("kvs").FindOne(context.TODO(), filter).Decode(&res)
					if err != nil {
						record := util.CmRecord{Name: posRecord.Value, Values: []string{}}
						_, err := db.Collection("kvs").InsertOne(context.TODO(), record)
						if err != nil {
							util.PrintErr(noStr, err)
						}
					}
				}

				/* Insert new value into final collection (if need be) */
				if insert == true && posEntry.Name != "Keys" {
					var res util.ValueEntry
					filterVal := bson.D{{Key: "name", Value: posEntry.Name}, // could be Keys
						{Key: "values", Value: posRecord.Value}}
					update := bson.D{{Key: "$push", Value: bson.D{
						{Key: "values", Value: posRecord.Value}}}}
					err := db.Collection("kvs").FindOne(context.TODO(), filterVal).Decode(&res)
					if err != nil {
						_, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
						if err != nil {
							util.PrintErr(noStr, err)
						}
					}
				}
			}
		}
	}
}

// check if the specified value entry in the value entry slice
func locateValEntry(arr []util.ValueEntry, val string) (util.ValueEntry, bool) {
	for _, valEntry := range arr {
		if valEntry.Value == val {
			return valEntry, true
		}
	}
	return util.ValueEntry{}, false
}
