package mgoToolkit

import (
	"context"

	"github.com/StevenZack/ghostman/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var clientPool = make(map[string]*mongo.Client)

func TakeClient(dsn string) (*mongo.Client, error) {
	if v, ok := clientPool[dsn]; ok {
		return v, nil
	}
	client, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	clientPool[dsn] = client
	return client, nil
}

func TakeDatabase(dsn string) (*mongo.Database, error) {
	info, e := connstring.Parse(dsn)
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	client, e := TakeClient(dsn)
	if e != nil {
		logx.Error(e)
		return nil, e
	}
	return client.Database(info.Database), nil
}
