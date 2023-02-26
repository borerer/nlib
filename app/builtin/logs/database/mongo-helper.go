package database

import "go.mongodb.org/mongo-driver/bson"

func FilterEquals(key string, val string) interface{} {
	return bson.M{
		key: val,
	}
}

var emptyFilter = bson.D{}
