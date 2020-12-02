package cryptoToolkit

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

func Sha1FromValues(vs map[string]interface{}) string {
	keys := []string{}
	for k := range vs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := new(strings.Builder)
	for i, key := range keys {
		buf.WriteString(key + "=" + fmt.Sprint(vs[key]))
		if i < len(keys)-1 {
			buf.WriteString("&")
		}
	}
	sha1 := sha1.New()
	io.Copy(sha1, strings.NewReader(buf.String()))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}
