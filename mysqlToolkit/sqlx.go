package mysqlToolkit

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SqlConn struct {
	db *sql.DB
}

func NewMySQL(dsn string) (*SqlConn, error) {
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		return nil, e
	}
	c := &SqlConn{
		db: db,
	}
	return c, nil
}

func (s *SqlConn) queryRow(v interface{}, query string, strict bool, args ...interface{}) error {
	t := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	if t.Kind().String() != "ptr" {
		return errors.New("value scan destination must be a pointer type")
	}
	t = t.Elem()
	value = value.Elem()

	rows, e := s.db.Query(query, args...)
	if e != nil {
		return e
	}

	columns, e := rows.Columns()
	if e != nil {
		return e
	}

	if strict && t.NumField() < len(columns) {
		return errors.New("value destination scan field num doesn't match")
	}

	scanArgs := make([]interface{}, len(columns))
	scanValues := make([]interface{}, len(columns))
	for i := range scanArgs {
		scanArgs[i] = &scanValues[i]
	}

	for rows.Next() {
		e := rows.Scan(scanArgs...)
		if e != nil {
			return e
		}
		break
	}

	fieldMap := GetDBTagMap(t)
	for i := 0; i < len(columns); i++ {
		fieldIndex, ok := fieldMap[columns[i]]
		if !ok {
			if strict {
				return errors.New("field '" + columns[i] + "' doesn't match")
			}
			continue
		}
		field := t.Field(fieldIndex)
		goValue, e := toGoValue(field, scanValues[i])
		if e != nil {
			return e
		}
		value.Field(fieldIndex).Set(reflect.ValueOf(goValue))
	}
	return nil
}

func (s *SqlConn) QueryRow(v interface{}, query string, args ...interface{}) error {
	return s.queryRow(v, query, true, args...)
}

func (s *SqlConn) QueryRowPartial(v interface{}, query string, args ...interface{}) error {
	return s.queryRow(v, query, false, args...)
}

func toGoValue(field reflect.StructField, v interface{}) (interface{}, error) {
	origin := string(v.([]uint8))
	println(origin, field.Type.Kind().String())
	switch field.Type.Kind() {
	case reflect.String:
		return origin, nil
		// int
	case reflect.Int:
		return strconv.Atoi(origin)
	case reflect.Int8:
		i, e := strconv.ParseInt(origin, 10, 64)
		return int8(i), e
	case reflect.Int16:
		i, e := strconv.ParseInt(origin, 10, 64)
		return int16(i), e
	case reflect.Int32:
		i, e := strconv.ParseInt(origin, 10, 64)
		return int32(i), e
	case reflect.Int64:
		i, e := strconv.ParseInt(origin, 10, 64)
		return i, e
		// uint
	case reflect.Uint:
		return strconv.Atoi(origin)
	case reflect.Uint8:
		i, e := strconv.ParseInt(origin, 10, 64)
		return uint8(i), e
	case reflect.Uint16:
		i, e := strconv.ParseInt(origin, 10, 64)
		return uint16(i), e
	case reflect.Uint32:
		i, e := strconv.ParseInt(origin, 10, 64)
		return uint32(i), e
	case reflect.Uint64:
		i, e := strconv.ParseInt(origin, 10, 64)
		return uint64(i), e
		// float
	case reflect.Float32:
		f, e := strconv.ParseFloat(origin, 64)
		return float32(f), e
	case reflect.Float64:
		f, e := strconv.ParseFloat(origin, 64)
		return f, e
	}
	// time
	if field.Type.Name() == "Time" {
		t, e := time.Parse(LayoutTime+" MST", origin+" CST")
		return t, e
	}
	return nil, errors.New("unsupported type " + field.Type.Kind().String())
}
