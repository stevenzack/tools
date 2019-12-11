package mysqlToolkit

import (
	"reflect"

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
