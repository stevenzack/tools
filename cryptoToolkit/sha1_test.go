package cryptoToolkit

import "testing"

func TestSha1FromValues(t *testing.T) {
	s := Sha1FromValues(map[string]interface{}{
		"two": 2,
		"one": 1,
	})
	s1 := `615f2daa674ed853197c0e4128757d885f5bd4d9`
	if s != s1 {
		t.Error("s is not s1 , but ", s)
		return
	}
}
