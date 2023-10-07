package mongo

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

// NewObjectId returns a new unique ObjectId.
func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}

func SafeToObjectId(s string) (bson.ObjectId, error) {
	if bson.IsObjectIdHex(s) {
		return bson.ObjectIdHex(s), nil
	} else {
		return bson.NewObjectId(), errors.New(fmt.Sprintf("INVALID_ID:%s", s))
	}
}

func SafeToObjectIdOrEmpty(s string) bson.ObjectId {
	if bson.IsObjectIdHex(s) {
		return bson.ObjectIdHex(s)
	} else {
		return ""
	}
}

func IsObjectIdHex(s string) bool {
	return bson.IsObjectIdHex(s)
}

func ToObjectIdHex(s string) bson.ObjectId {
	return bson.ObjectIdHex(s)
}

// m bson.M
// s SomeStruct
// BsonM2Struct(m, &s)
func BsonM2Struct(m bson.M, s interface{}) (status bool) {
	if bsonBytes, err := bson.Marshal(m); err == nil {
		bson.Unmarshal(bsonBytes, s)
		status = true
	} else {
		status = false
	}
	return
}

func FuzzyQuery(value interface{}) bson.M {
	return bson.M{"$regex": bson.RegEx{value.(string), "i"}}
}

func SizeQuery(size int) bson.M {
	return bson.M{"$size": size}
}

func InQuery(valueList interface{}) bson.M {
	return bson.M{"$in": valueList}
}

func NinQuery(valueList interface{}) bson.M {
	return bson.M{"$nin": valueList}
}

func OrQuery(qs []bson.M) bson.M {
	return bson.M{"$or": qs}
}

func AndQuery(qs []bson.M) bson.M {
	return bson.M{"$and": qs}
}

func SetOperator(value interface{}) bson.M {
	return bson.M{"$set": value}
}

func IncOperator(valueList interface{}) bson.M {
	return bson.M{"$inc": valueList}
}

func NotQuery(value interface{}) bson.M {
	return bson.M{"$ne": value}
}

func NExistQuery() bson.M {
	return bson.M{"$exists": false}
}
