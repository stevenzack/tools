package mongodbx

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"

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
	groups := make(map[string]mongo.IndexModel)
	for key, v := range indexes {
		vs, e := url.ParseQuery(strings.ReplaceAll(v, ",", "&"))
		if e != nil {
			return errors.New("field '" + key + "', invalid value format:" + v)
		}
		groupseq := 1
		unique := false
		group := ""
		single := 0
		for k := range vs {
			switch k {
			case "groupseq":
				sequence := vs.Get(k)
				if sequence != "" {
					groupseq, e = strconv.Atoi(sequence)
					if e != nil {
						return errors.New("field '" + key + "', invalid groupseq format:" + v)
					}
					if groupseq != -1 {
						groupseq = 1
					}
				}
			case "unique":
				unique = vs.Get("unique") == "true"
			case "group":
				group = vs.Get(k)
			case "single":
				sequence := vs.Get(k)
				if sequence != "" {
					single, e = strconv.Atoi(sequence)
					if e != nil {
						return errors.New("field '" + key + "', invalid single format:" + v)
					}
					if single != -1 {
						single = 1
					} else {
						single = -1
					}
				}
			default:
				return errors.New("field '" + key + "', unsupported key:" + k)
			}
		}

		if group == "" {
			single = 1
		}

		//single index
		if single != 0 {
			imodel := mongo.IndexModel{
				Keys: bson.D{
					{
						Key:   key,
						Value: single,
					},
				},
				Options: options.Index(),
			}
			if unique {
				imodel.Options.SetUnique(unique)
			}
			imodels = append(imodels, imodel)
		}
		//group index
		if group != "" {
			imodel, ok := groups[group]
			if !ok {
				imodel = mongo.IndexModel{
					Keys: bson.D{
						{
							Key:   key,
							Value: groupseq,
						},
					},
					Options: options.Index(),
				}
				if strings.HasPrefix(group, "unique") {
					imodel.Options.SetUnique(unique)
				}
				groups[group] = imodel
				continue
			}
			imodel.Keys = append(imodel.Keys.(bson.D), bson.E{
				Key:   key,
				Value: groupseq,
			})
			groups[group] = imodel
		}
	}

	//add group indexes
	for _, v := range groups {
		imodels = append(imodels, v)
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
