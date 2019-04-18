package fileToolkit

import "testing"

func Test_a(t *testing.T) {
	dir := GetPathOfGoPkg("/Users/stevenzacker/go/src/github.com/StevenZack/gengo/example/data", "../model")
	t.Log(GetGoPkgFromDirPath(dir))
}
