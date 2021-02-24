package strToolkit

import "testing"

func TestToVersionCodeFloat(t *testing.T) {
	f, e := ToVersionCodeFloat("1.2.3.4.5")
	if e != nil {
		t.Error(e)
		return
	}
	if f != 10203.0405 {
		t.Error("f is not 10203.0405 , but ", f)
		return
	}
}
