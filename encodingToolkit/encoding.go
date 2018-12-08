package encodingToolkit

import (
	"net/url"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func EncodingGBK(src string) string {
	if runtime.GOOS != "windows" {
		return src
	}
	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), simplifiedchinese.GBK.NewEncoder()))
	return string(data)
}
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
func UrlEncode(s string)string{
	return url.QueryEscape(s)
}
func UrlDecode(s string)(string,error){
	return url.QueryUnescape(s)
}