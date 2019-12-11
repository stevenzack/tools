package mysqlToolkit

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func ParseMySQLDatabase(dsn string) (string, string, error) {
	info, e := mysql.ParseDSN(dsn)
	if e != nil {
		return "", "", e
	}
	db := info.DBName
	if db == "" {
		return "", "", errors.New("连接串没写database")
	}
	info.Loc = time.Local
	info.ParseTime = true
	return info.FormatDSN(), db, nil
}

func TableExists(conn *SqlConn, table string) (bool, error) {
	_, e := conn.Exec(`describe ` + table)
	if e != nil {
		if strings.Contains(e.Error(), `Table`) && strings.Contains(e.Error(), `doesn't exist`) {
			return false, nil
		}
		return false, e
	}
	return true, nil
}

func CreateTableIfNotExists(conn *SqlConn, table, sql string) (bool, error) {
	exists, e := TableExists(conn, table)
	if e != nil {
		return false, e
	}
	if exists {
		return false, nil
	}
	_, e = conn.Exec(sql)
	if e != nil {
		return false, e
	}
	return true, nil
}

func CreateIndexes(conn *SqlConn, db, table string, indexes []string) error {
	for _, index := range indexes {
		_, e := conn.Exec(`create index ` + index + `_idx on ` + db + `.` + table + ` (` + index + `)`)
		if e != nil {
			return e
		}
	}
	return nil
}

func toGoValue(field reflect.StructField, mysqlValue string) (interface{}, error) {
	origin := mysqlValue
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
