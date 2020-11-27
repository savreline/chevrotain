package main

import (
	"context"

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
