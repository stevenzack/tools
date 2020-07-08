package mgoToolkit

import (
	"context"
	"errors"
	"log"
	"reflect"

	"github.com/StevenZack/tools/strToolkit"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseModel struct {
	DataSourceName string //data source name
	CollectionName string // collection name
	Type           reflect.Type
	Data           interface{}
	Collection     *mongo.Collection
}

func NewBaseModel(dsn string, data interface{}) (*BaseModel, error) {
	model := &BaseModel{DataSourceName: dsn}
	e := model.initData(data)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	model.Collection, e = model.takeCollection()
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return model, nil
}

func (b *BaseModel) initData(data interface{}) error {
	t := reflect.TypeOf(data)
	b.Type = t
	b.CollectionName = strToolkit.LowerFirst(t.Name())

	if t.Kind().String() == "ptr" {
		return errors.New("data必须是非指针类型")
	}

	if t.Field(0).Type.Name() != "ObjectID" {
		return errors.New(t.Name() + "类型的第一个字段不是primitive.ObjectID类型")
	}

	indexes := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bson, ok := field.Tag.Lookup("bson")
		if !ok {
			return errors.New(t.Name() + "类型的" + field.Name + "字段没有加bson的tag")
		}

		if _, ok := field.Tag.Lookup("index"); ok || bson == "createTime" {
			indexes[bson] = 1
		}
	}

	db, e := TakeDatabase(b.DataSourceName)
	if e != nil {
		log.Println(e)
		return e
	}
	e = CreateIndexIfNotExists(db, b.CollectionName, indexes)
	if e != nil {
		log.Println(e)
		return e
	}
	return nil
}

func (b *BaseModel) takeCollection() (*mongo.Collection, error) {
	db, e := TakeDatabase(b.DataSourceName)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	return db.Collection(b.CollectionName), nil
}

func (b *BaseModel) Insert(v interface{}) error {
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	if t.Kind().String() == "ptr" {
		t = t.Elem()
		value = value.Elem()
	}
	if t.Name() != b.Type.Name() {
		return errors.New("插入的数据不是" + b.Type.Name() + "类型")
	}

	objValue := value.Field(0).Interface()
	objId, ok := objValue.(primitive.ObjectID)
	if !ok {
		return errors.New("插入的数据的第一个值不是primitive.ObjectID类型")
	}
	if objId == primitive.NilObjectID {
		value.Field(0).Set(reflect.ValueOf(primitive.NewObjectID()))
	}

	coll, e := b.takeCollection()
	if e != nil {
		log.Println(e)
		return e
	}

	_, e = coll.InsertOne(context.TODO(), v)
	if e != nil {
		log.Println(e)
		return e
	}
	return nil
}

func (b *BaseModel) Find(id string) (interface{}, error) {
	objId, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	coll, e := b.takeCollection()
	if e != nil {
		log.Println(e)
		return nil, e
	}

	v := reflect.New(b.Type)
	e = coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(v.Interface())
	if e != nil {
		log.Println(e)
		return nil, e
	}

	return v.Interface(), nil
}

func (b *BaseModel) Update(id string, updater bson.M) (int64, error) {
	objId, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		log.Println(e)
		return 0, e
	}

	coll, e := b.takeCollection()
	if e != nil {
		log.Println(e)
		return 0, e
	}

	l, e := coll.UpdateOne(context.TODO(), bson.M{"_id": objId}, bson.M{
		"$set": updater,
	})
	if e != nil {
		log.Println(e)
		return 0, e
	}
	return l.ModifiedCount, nil
}
