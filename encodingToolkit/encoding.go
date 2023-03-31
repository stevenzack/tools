package encodingToolkit

import (
	"encoding/json"
	"net/url"
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
func UrlEncode(s string) string {
	return url.QueryEscape(s)
}
func UrlDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}
