package strToolkit

import (
	"testing"

	"github.com/StevenZack/tools/numToolkit"
)

func TestBaseDecode(t *testing.T) {
	i := numToolkit.Rand63n(100000)
	s := BaseEncode(i)
	t.Log(s)
	i2, e := BaseDecode(s)
	if e != nil {
		t.Error(e)
		return
	}
	if i != i2 {
		t.Error("i is not i2 , but ", i)
		return
	}
}
