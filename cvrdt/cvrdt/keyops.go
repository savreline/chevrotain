package main

import (
	"context"

	"../../util"
	"github.com/savreline/GoVector/govec"
)

// InsertKey inserts the given key with an empty array for values
func (t *RPCExt) InsertKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, posCollection, nil)
	return nil
}

// RemoveKey removes the given key
func (t *RPCExt) RemoveKey(args *util.RPCExtArgs, reply *int) error {
	InsertKeyLocal(args.Key, negCollection, nil)
	return nil
}

// InsertKeyLocal inserts the key in either positive collection (add) or negative collection (remove)
func InsertKeyLocal(key string, collection string, record *util.CvRecord) {
	/* In no ready to go record is supplied, tick the clock and make one */
	if record == nil {
		logger.LogLocalEvent("Inserting Key "+key, govec.GetDefaultLogOptions())
		record = &util.CvRecord{Name: key, Timestamp: logger.GetCurrentVC(), Values: []util.ValueEntry{}}
	} else {
		if len(record.Values) == 0 {
			record.Values = []util.ValueEntry{}
		}
	}

	/* Insert the Key */
	_, err := db.Collection(collection).InsertOne(context.TODO(), record)
	if err != nil {
		util.PrintErr(noStr, err)
	}
	util.PrintMsg(noStr, "Inserted Key "+key)
}
