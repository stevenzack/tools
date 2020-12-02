package cryptoToolkit

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
)

func Sha1FromValues(vs map[string]interface{}) string {
	keys := []string{}
	for k := range vs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	values := url.Values{}
	for _, key := range keys {
		values.Add(key, fmt.Sprint(vs[key]))
	}
	str := values.Encode()
	sha1 := sha1.New()
	io.Copy(sha1, strings.NewReader(str))
	return fmt.Sprintf("%x", sha1.Sum(nil))
}
