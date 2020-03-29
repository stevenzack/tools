package mysqlToolkit

import (
	"database/sql"
	"errors"
	"reflect"

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
	if t.Kind() != reflect.Ptr {
		return errors.New("scan destination must be a pointer type")
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
		return errors.New("scan destination's field num doesn't match")
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
		goValue, e := toGoValue(field, string(scanValues[i].([]uint8)))
		if e != nil {
			return e
		}
		value.Field(fieldIndex).Set(reflect.ValueOf(goValue))
	}
	return nil
}

func (s *SqlConn) queryRows(v interface{}, query string, strict bool, args ...interface{}) error {
	t := reflect.TypeOf(v)
	sliceValue := reflect.ValueOf(v)
	if t.Kind() != reflect.Ptr {
		return errors.New(" scan destination must be a pointer type")
	}
	t = t.Elem()
	sliceValue = sliceValue.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New("scan destination must be a pointer of slice")
	}
	t = t.Elem()
	ptrMode := false
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		ptrMode = true
	}

	rows, e := s.db.Query(query, args...)
	if e != nil {
		return e
	}

	columns, e := rows.Columns()
	if e != nil {
		return e
	}

	if strict && t.NumField() < len(columns) {
		return errors.New("scan destination's field num doesn't match")
	}
	scanArgs := make([]interface{}, len(columns))
	scanValues := make([]interface{}, len(columns))
	for i := range scanArgs {
		scanArgs[i] = &scanValues[i]
	}

	fieldMap := GetDBTagMap(t)
	for rows.Next() {
		e := rows.Scan(scanArgs...)
		if e != nil {
			return e
		}
		value := reflect.New(t)

		for i := 0; i < len(columns); i++ {
			fieldIndex, ok := fieldMap[columns[i]]
			if !ok {
				if strict {
					return errors.New("field '" + columns[i] + "' doesn't match")
				}
				continue
			}
			field := t.Field(fieldIndex)
			goValue, e := toGoValue(field, string(scanValues[i].([]uint8)))
			if e != nil {
				return e
			}
			value.Elem().Field(fieldIndex).Set(reflect.ValueOf(goValue))
		}

		if !ptrMode {
			value = value.Elem()
		}
		sliceValue.Set(reflect.Append(sliceValue, value))
	}

	return nil
}

func (s *SqlConn) QueryRow(v interface{}, query string, args ...interface{}) error {
	return s.queryRow(v, query, true, args...)
}

func (s *SqlConn) QueryRowPartial(v interface{}, query string, args ...interface{}) error {
	return s.queryRow(v, query, false, args...)
}

func (s *SqlConn) QueryRows(vs interface{}, query string, args ...interface{}) error {
	return s.queryRows(vs, query, true, args...)
}

func (s *SqlConn) QueryRowsPartial(vs interface{}, query string, args ...interface{}) error {
	return s.queryRows(vs, query, false, args...)
}

func (s *SqlConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args)
}

type Column struct {
	Field string `db:"Field"`
}

func DescribeTable(conn *SqlConn, table string) ([]*Column, error) {
	vs := []*Column{}
	e := conn.QueryRows(&vs, `describe `+table)
	return vs, e
}
