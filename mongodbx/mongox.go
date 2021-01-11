package mongodbx

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func CreateIndex(coll *mongo.Collection, indexes map[string]string) error {
	imodels := []mongo.IndexModel{}
	for key, v := range indexes {
		vs, e := url.ParseQuery(v)
		if e != nil {
			return errors.New("field '" + key + "', invalid value format:" + v)
		}
		seq := -1
		options := options.Index()
		for k := range vs {
			switch k {
			case "seq":
				sequence := vs.Get(k)
				if sequence != "" {
					seq, e = strconv.Atoi(sequence)
					if e != nil {
						return errors.New("field '" + key + "', invalid seq format:" + v)
					}
					if seq != -1 {
						seq = 1
					}
				}
			case "uniq":
				uniq := vs.Get("uniq")
				if uniq != "" {
					options.SetUnique(uniq == "true")
				}
			default:
				return errors.New("field '" + key + "', unsupported key:" + k)
			}
		}

		imodel := mongo.IndexModel{
			Keys: bson.M{
				key: seq,
			},
			Options: options,
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
func CreateIndexIfNotExists(db *mongo.Database, collname string, indexes map[string]string) (bool, error) {
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
