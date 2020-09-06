package mgoToolkit


import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// DialCollection DialCollection
func DialCollection(db *mongo.Database, name string) *mongo.Collection {
	coll := db.Collection(name)
	return coll
}

// CollectionExists check if collection exists
func CollectionExists(db *mongo.Database, coll string) (bool, error) {
	names, e := db.ListCollectionNames(context.TODO(), bson.M{})
	if e != nil {
		return false, e
	}

	for _, name := range names {
		if name == coll {
			return true, nil
		}
	}
	return false, nil
}

// CreateIndex creates indexes for coll
func CreateIndex(coll *mongo.Collection, indexes map[string]int) error {
	imodels := []mongo.IndexModel{}
	for k, v := range indexes {
		sequence := bsonx.Int32(1)
		if v < 0 {
			sequence = bsonx.Int32(-1)
		}
		imodel := mongo.IndexModel{
			Keys: bsonx.Doc{bsonx.Elem{Key: k, Value: sequence}},
		}
		imodels = append(imodels, imodel)
	}
	if len(imodels) == 0 {
		return nil
	}

	_, e := coll.Indexes().CreateMany(context.TODO(), imodels)
	return e
}

// CreateIndexIfNotExists create indexes if collection doesn't exists
func CreateIndexIfNotExists(db *mongo.Database, collname string, indexes map[string]int) (bool, error) {
	b, e := CollectionExists(db, collname)
	if e != nil {
		return false, e
	}
	if b {
		return false, nil
	}

	coll := DialCollection(db, collname)

	return true, CreateIndex(coll, indexes)
}
