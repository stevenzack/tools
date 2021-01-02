package strToolkit

import (
	"testing"
)

func TestSubBefore(t *testing.T) {
	s := SubBefore("28341/af1", "/a", "")
	if s != `28341` {
		t.Error("s is not `28341` , but ", s)
		return
	}
}

func TestSubBeforeLast(t *testing.T) {
	s := SubBeforeLast("12/3/313", "/3", "")
	if s != `12/3` {
		t.Error("s is not `12/3` , but ", s)
		return
	}
}

func TestSubAfter(t *testing.T) {
	s := SubAfter("123/a38/a00", "/a", "")
	if s != `38/a00` {
		t.Error("s is not `38/a00` , but ", s)
		return
	}
}

func TestSubAfterLast(t *testing.T) {
	s := SubAfterLast("123/a38/a00", "/a", "")
	if s != `00` {
		t.Error("s is not `00` , but ", s)
		return
	}
}
