package mysqlToolkit

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func GetDBTagMap(t reflect.Type) map[string]int {
	m := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		m[GetDBTag(field)] = i
	}
	return m
}

func GetDBTag(field reflect.StructField) string {
	db, ok := field.Tag.Lookup("db")
	if !ok {
		return strToolkit.ToSnakeCase(field.Name)
	}
	return db
}

func GoTypeToSQLType(t reflect.Type, db string, length int) (string, bool, error) {
	modifiable := true
	switch t.Name() {
	case "string":
		if t.Name() == "string" {
			if length == 0 {
				return "text", modifiable, nil
			}
			return "varchar(" + strconv.Itoa(length) + ")", modifiable, nil
		}
	case "Time":
		if db == "create_time" || db == "update_time" {
			return "timestamp NULL default CURRENT_TIMESTAMP", false, nil
		}
		return "timestamp NULL", modifiable, nil
	}

	if strings.Contains(t.Name(), "uint") {
		if length != 0 {
			return "int(" + strconv.Itoa(length) + ") unsigned", modifiable, nil
		}
		return "int unsigned", modifiable, nil
	}

	if strings.Contains(t.Name(), "int") {
		if length != 0 {
			return "int(" + strconv.Itoa(length) + ")", modifiable, nil
		}
		return "int", modifiable, nil
	}
	if strings.Contains(t.Name(), "float") {
		if length != 0 {
			return "float(" + strconv.Itoa(length) + ")", modifiable, nil
		}
		return "float", modifiable, nil
	}
	return "", modifiable, errors.New("unsupported Go type:" + t.Name())
}
