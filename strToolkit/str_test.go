package strToolkit

import "testing"

func Test_compare(t *testing.T) {
	t.Log(CompareVersionLeftHigher("1.10.0", "1.0.0"))
}
