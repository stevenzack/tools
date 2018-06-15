package strToolkit

import (
	"encoding/json"
)

func JsonObj(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		return "{}"
	}
	return string(b)
}
func JsonArray(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		return "[]"
	}
	return string(b)
}
