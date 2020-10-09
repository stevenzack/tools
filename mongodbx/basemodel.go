package mongodbx

import (
	"context"
	"errors"
	"reflect"

	"github.com/iancoleman/strcase"
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
	model, _, e := NewBaseModelWithCreated(dsn, data)
	return model, e
}

func NewBaseModelWithCreated(dsn string, data interface{}) (*BaseModel, bool, error) {
	model := &BaseModel{DataSourceName: dsn}
	created, e := model.initData(data)
	if e != nil {
		return nil, false, e
	}
	model.Collection, e = model.takeCollection()
	if e != nil {
		return nil, false, e
	}

	return model, created, nil
}

func (b *BaseModel) initData(data interface{}) (bool, error) {
	t := reflect.TypeOf(data)
	b.Type = t
	b.CollectionName = strcase.ToLowerCamel(t.Name())

	if t.Kind().String() == "ptr" {
		return false, errors.New("data必须是非指针类型")
	}

	indexes := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if i == 0 {
			if field.Type.String() != "primitive.ObjectID" {
				return false, errors.New(t.Name() + "类型的" + field.Name + "字段不是primitive.ObjectID类型")
			}
		}
		bson, ok := field.Tag.Lookup("bson")
		if !ok {
			return false, errors.New(t.Name() + "类型的" + field.Name + "字段没有加bson的tag")
		}
		if i == 0 {
			if bson != "_id,omitempty" {
				return false, errors.New(t.Name() + "类型的" + field.Name + "字段tag不是 _id,omitempty")
			}
		}

		if _, ok := field.Tag.Lookup("index"); ok || bson == "createTime" {
			indexes[bson] = 1
		}
	}

	db, e := TakeDatabase(b.DataSourceName)
	if e != nil {
		return false, e
	}
	created, e := CreateIndexIfNotExists(db, b.CollectionName, indexes)
	if e != nil {
		return false, e
	}
	return created, nil
}

func (b *BaseModel) takeCollection() (*mongo.Collection, error) {
	db, e := TakeDatabase(b.DataSourceName)
	if e != nil {
		return nil, e
	}
	return db.Collection(b.CollectionName), nil
}

func (b *BaseModel) Insert(v interface{}) (string, error) {
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	if t.Kind().String() == "ptr" {
		t = t.Elem()
		value = value.Elem()
	}
	if t.Name() != b.Type.Name() {
		return "", errors.New("插入的数据不是" + b.Type.Name() + "类型")
	}

	coll, e := b.takeCollection()
	if e != nil {
		return "", e
	}

	result, e := coll.InsertOne(context.TODO(), v)
	if e != nil {
		return "", e
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (b *BaseModel) Find(id string) (interface{}, error) {
	coll, e := b.takeCollection()
	if e != nil {
		return nil, e
	}

	obj, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return nil, errors.New(e.Error() + ":" + id)
	}
	v := reflect.New(b.Type)
	e = coll.FindOne(context.TODO(), bson.M{"_id": obj}).Decode(v.Interface())
	if e != nil {
		return nil, e
	}

	return v.Interface(), nil
}

func (b *BaseModel) FindWhere(where bson.M) (interface{}, error) {
	coll, e := b.takeCollection()
	if e != nil {
		return nil, e
	}
	v := reflect.New(b.Type)
	e = coll.FindOne(context.TODO(), where).Decode(v.Interface())
	if e != nil {
		return nil, e
	}

	return v.Interface(), nil
}

func (b *BaseModel) QueryWhere(where bson.M) (interface{}, error) {
	coll, e := b.takeCollection()
	if e != nil {
		return nil, e
	}
	vs := reflect.New(reflect.SliceOf(reflect.PtrTo(b.Type)))
	cursor, e := coll.Find(context.TODO(), where)
	if e != nil {
		return nil, e
	}

	defer cursor.Close(context.TODO())
	e = cursor.All(context.TODO(), vs.Interface())
	if e != nil {
		return nil, e
	}
	return vs.Elem().Interface(), nil
}

func (b *BaseModel) UpdateSet(id string, updater bson.M) (int64, error) {
	coll, e := b.takeCollection()
	if e != nil {
		return 0, e
	}
	obj, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return 0, nil
	}

	l, e := coll.UpdateOne(context.TODO(), bson.M{"_id": obj}, bson.M{
		"$set": updater,
	})
	if e != nil {
		return 0, e
	}
	return l.ModifiedCount, nil
}

func (b *BaseModel) Update(id string, updator bson.M) (int64, error) {
	coll, e := b.takeCollection()
	if e != nil {
		return 0, e
	}

	obj, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return 0, e
	}

	l, e := coll.UpdateOne(context.TODO(), bson.M{"_id": obj}, updator)
	if e != nil {
		return 0, e
	}

	return l.ModifiedCount, nil
}

func (b *BaseModel) UpdateWhere(where, updator bson.M) (int64, error) {
	r, e := b.Collection.UpdateMany(context.TODO(), where, updator)
	if e != nil {
		return 0, e
	}
	return r.ModifiedCount, nil
}

// Clear clear collection
func (b *BaseModel) Clear() error {
	_, e := b.Collection.DeleteMany(context.TODO(), bson.M{})
	return e
}

func (b *BaseModel) Delete(id string) (int64, error) {
	obj, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return 0, e
	}
	r, e := b.Collection.DeleteOne(context.TODO(), bson.M{
		"_id": obj,
	})
	if e != nil {
		return 0, e
	}
	return r.DeletedCount, nil
}

func (b *BaseModel) DeleteWhere(where bson.M) (int64, error) {
	r, e := b.Collection.DeleteMany(context.TODO(), where)
	if e != nil {
		return 0, e
	}
	return r.DeletedCount, nil
}