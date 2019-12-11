package mysqlToolkit

import (
	"reflect"
	"strconv"
)

func GetLengthTag(field reflect.StructField) (int, error) {
	length, ok := field.Tag.Lookup("length")
	if !ok {
		return 0, nil
	}
	return strconv.Atoi(length)
}
