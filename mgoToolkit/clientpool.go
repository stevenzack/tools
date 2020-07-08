package mgoToolkit

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var clientPool = sync.Map{}

func TakeClient(dsn string) (*mongo.Client, error) {
	v, ok := clientPool.Load(dsn)
	if !ok {
		mgo, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
		if e != nil {
			log.Println(e)
			return nil, e
		}
		clientPool.Store(dsn, mgo)
		return mgo, nil
	}
	return v.(*mongo.Client), nil
}

func TakeDatabase(dsn string) (*mongo.Database, error) {
	info, e := connstring.Parse(dsn)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	client, e := TakeClient(dsn)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return client.Database(info.Database), nil
}
