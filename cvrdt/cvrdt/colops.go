package main

import (
	"context"
	"fmt"

	"../../util"
	"go.mongodb.org/mongo-driver/bson"
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
	/* Merge the global keys entry */
	var keysPosRecord, keysNegRecord util.CvRecord
	filter := bson.D{{Key: "name", Value: "Keys"}}
	err := db.Collection(posCollection).FindOne(context.TODO(), filter).Decode(&keysPosRecord)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	err = db.Collection(negCollection).FindOne(context.TODO(), filter).Decode(&keysNegRecord)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	mergeArray(keysPosRecord.Values, keysNegRecord.Values, "")

	/* Merge all other entries */
	// cursor, err := db.Collection(posCollection).Find(context.TODO(), bson.D{})
	// if err != nil {
	// 	fmt.Println("error 1")
	// 	util.PrintErr(noStr, err)
	// }
	// for cursor.Next(context.TODO()) {
	// 	/* Get next item */
	// 	var posRecord, negRecord util.CvRecord
	// 	if err := cursor.Decode(&posRecord); err != nil {
	// 		fmt.Println("error 2")
	// 		util.PrintErr(noStr, err)
	// 	}

	// 	/* Look in the negative collection */
	// 	filter = bson.D{{Key: "name", Value: posRecord.Name}}
	// 	err = db.Collection(negCollection).FindOne(context.TODO(), filter).Decode(&negRecord)
	// 	if err != nil {
	// 		fmt.Println("error 3")
	// 		util.PrintErr(noStr, err)
	// 	}
	// 	mergeArray(posRecord.Values, negRecord.Values, posRecord.Name)
	// }
}

// merge specific arrays
func mergeArray(posArray []util.ValueEntry, negArray []util.ValueEntry, key string) {
	for _, posRecord := range posArray {
		/* Since item been in posCollection, default action is to insert it */
		var insert = true
		filter := bson.D{{Key: "name", Value: key}}

		/* Try to locate the item in the negative collection */
		negRecord, found := locateValEntry(negArray, posRecord.Value)
		if found == true {
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

			/* Delete positive record */
			update := bson.D{{Key: "$pull", Value: bson.D{
				{Key: "values", Value: posRecord}}}}
			_, err := db.Collection(posCollection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				fmt.Println("error 4")
				util.PrintErr(noStr, err)
			}

			/* Delete negative record */
			update = bson.D{{Key: "$pull", Value: bson.D{
				{Key: "values", Value: negRecord}}}}
			_, err = db.Collection(negCollection).UpdateOne(context.TODO(), filter, update)
			if err != nil {
				fmt.Println("error 5")
				util.PrintErr(noStr, err)
			}

			if key == "" {
				keyFilter := bson.D{{Key: "name", Value: posRecord.Value}}

				/* Delete key in positive collection */
				_, err := db.Collection(posCollection).DeleteOne(context.TODO(), keyFilter)
				if err != nil {
					fmt.Println("error 6")
					util.PrintErr(noStr, err)
				}

				/* Delete key in negative collection */
				_, err = db.Collection(negCollection).DeleteOne(context.TODO(), keyFilter)
				if err != nil {
					fmt.Println("error 7")
					util.PrintErr(noStr, err)
				}
			}
		}

		/* Insert new key, if need be */
		if insert == true && key == "" {
			record := util.CmRecord{Name: posRecord.Value, Values: []string{}}
			_, err := db.Collection("kvs").InsertOne(context.TODO(), record)
			if err != nil {
				fmt.Println("error 8")
				util.PrintErr(noStr, err)
			}
		}

		/* Insert new value, if need be */
		if insert == true && key != "" {
			filter = bson.D{{Key: "name", Value: key}}
			update := bson.D{{Key: "$push", Value: bson.D{
				{Key: "values", Value: posRecord.Value}}}}
			_, err := db.Collection("kvs").UpdateOne(context.TODO(), filter, update)
			if err != nil {
				fmt.Println("error 9")
				util.PrintErr(noStr, err)
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
