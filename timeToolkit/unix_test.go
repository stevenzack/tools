package timeToolkit

import (
	"testing"
)

func TestParseUnix(t *testing.T) {
	v, e := ParseUnix("1625492607")
	if e != nil {
		t.Error(e)
		return
	}
	s := v.Format(LAYOUT_DATE)
	if s != `2021-07-05` {
		t.Error("s is not `2021-07-05` , but ", s)
		return
	}

	v, e = ParseUnix("1625492607000")
	if e != nil {
		t.Error(e)
		return
	}
	s = v.Format(LAYOUT_DATE)
	if s != `2021-07-05` {
		t.Error("s is not `2021-07-05` , but ", s)
		return
	}

	v, e = ParseUnix("162549260700")
	if e != nil {
		t.Error(e)
		return
	}
	s = v.Format(LAYOUT_DATE)
	if s != `2021-07-05` {
		t.Error("s is not `2021-07-05` , but ", s)
		return
	}
}
