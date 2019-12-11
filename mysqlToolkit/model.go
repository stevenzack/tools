package mysqlToolkit

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

type (
	BaseMySQLModel struct {
		Conn      *SqlConn
		AppName   string
		Database  string
		TableName string
		Type      reflect.Type

		DBs              []string
		ModifiableDBs    []string
		ModifiableFields []string

		sqlGenerated string
		sqlInsert    string
	}
	Count struct {
		Count int64 `db:"count"`
	}
)

// NewBaseMySQLModel 新建基础Model
func NewBaseMySQLModel(appName string, dsn string, data interface{}) (*BaseMySQLModel, error) {
	b := &BaseMySQLModel{
		AppName: appName,
	}
	var e error
	dsn, b.Database, e = ParseMySQLDatabase(dsn)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	b.Conn, e = NewMySQL(dsn)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	e = b.initData(data)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	b.generateSQLInsert()
	return b, nil
}

// initData 如果表不存在，会自动建表，建索引
func (b *BaseMySQLModel) initData(data interface{}) error {
	b.Type = reflect.TypeOf(data)
	if b.Type.Kind().String() == "ptr" {
		b.Type = b.Type.Elem()
	}
	b.TableName = strToolkit.ToSnakeCase(b.AppName) + "_" + strToolkit.ToSnakeCase(b.Type.Name())

	indexes := []string{}
	b.sqlGenerated = `create table ` + b.Database + `.` + b.TableName + "(\n"
	for i := 0; i < b.Type.NumField(); i++ {
		field := b.Type.Field(i)
		// tag check
		db, ok := field.Tag.Lookup("db")
		if !ok {
			return errors.New(b.Type.Name() + "类型的" + field.Name + "字段没有写'db' Tag")
		}
		if db != strToolkit.ToSnakeCase(field.Name) {
			return errors.New(b.Type.Name() + "类型的'db'Tag格式不是标准的SnakeCase")
		}
		comment, ok := field.Tag.Lookup("comment")
		if !ok {
			return errors.New(b.Type.Name() + "类型的" + field.Name + "字段没写comment")
		}
		length, e := GetLengthTag(field)
		if e != nil {
			return errors.New(b.Type.Name() + "类型的" + field.Name + "字段的'length' Tag格式不正确")
		}

		// collect info
		b.DBs = append(b.DBs, db)
		if _, ok := field.Tag.Lookup("index"); ok {
			indexes = append(indexes, db)
		}

		// sql generating
		sqlType, modifiable, e := GoTypeToSQLType(field.Type, db, length)
		if e != nil {
			log.Println(e)
			return e
		}
		b.sqlGenerated += db + ` ` + sqlType + ` `
		if i == 0 {
			if strings.Contains(sqlType, "int") {
				b.sqlGenerated += `auto_increment `
				modifiable = false
			}
		} else if !strings.Contains(sqlType, "timestamp") {
			b.sqlGenerated += `not null `
		}
		if db == `update_time` {
			b.sqlGenerated += ` on update CURRENT_TIMESTAMP `
		}
		b.sqlGenerated += ` comment '` + comment + `' `
		if i == 0 {
			b.sqlGenerated += ` primary key `
		}
		b.sqlGenerated += ",\n"

		if modifiable {
			b.ModifiableDBs = append(b.ModifiableDBs, db)
			b.ModifiableFields = append(b.ModifiableFields, field.Name)
		}
	}
	// sql
	b.sqlGenerated = strings.TrimSuffix(b.sqlGenerated, ",\n")
	b.sqlGenerated += "\n)"

	created, e := CreateTableIfNotExists(b.Conn, b.TableName, b.sqlGenerated)
	if e != nil {
		log.Println(e)
		return e
	}
	if created {
		e := CreateIndexes(b.Conn, b.Database, b.TableName, indexes)
		if e != nil {
			log.Println(e)
			return e
		}
	}

	return nil
}

func (b *BaseMySQLModel) generateSQLInsert() {
	b.sqlInsert = `insert into ` + b.TableName + ` (` + strings.Join(b.ModifiableDBs, ",") + `) values(
		` + strings.Join(strToolkit.SlicifyStr("?", len(b.ModifiableDBs)), ",") + `
		)`
}

func (b *BaseMySQLModel) Insert(data interface{}) error {
	args := []interface{}{}
	value := reflect.ValueOf(data)
	if value.Kind().String() == "ptr" {
		value = value.Elem()
	}

	for _, name := range b.ModifiableFields {
		args = append(args, value.FieldByName(name).Interface())
	}

	_, e := b.Conn.Exec(b.sqlInsert, args...)
	return e
}

func (b *BaseMySQLModel) FindBy(field string, fieldValue interface{}) (interface{}, error) {
	if !strToolkit.SliceContains(b.DBs, field) {
		return nil, errors.New(`field ` + field + ` doesn't exists`)
	}

	query := `select ` + strings.Join(b.DBs, ",") + ` from ` + b.TableName + ` where ` + field + `=?`
	v := reflect.New(b.Type).Interface()
	e := b.Conn.QueryRow(v, query, fieldValue)
	if e != nil {
		return nil, e
	}
	return v, nil
}

func (b *BaseMySQLModel) Count(where string, args ...interface{}) (int64, error) {
	c := Count{}
	if where != "" {
		where = ` where ` + where
	}
	e := b.Conn.QueryRow(&c, `select count() as count from `+b.TableName+where, args...)
	if e != nil {
		log.Println(e)
		return 0, e
	}
	return c.Count, nil
}
