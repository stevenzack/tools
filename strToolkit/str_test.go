package strToolkit

import (
	"testing"
)

func Test_compare(t *testing.T) {
	t.Log(ToSnakeCase("UsernameAge"))
}

func TestSubBefore(t *testing.T) {
	s := SubBefore("asd@94@yes", "@")
	if s != `asd` {
		t.Error("s is not `asd` , but ", s)
		return
	}
}

func TestSubBeforeLast(t *testing.T) {
	s := SubBeforeLast("asd@87@yes", "@")
	if s != `asd@87` {
		t.Error("s is not `asd@87` , but ", s)
		return
	}
}

func TestSubAfter(t *testing.T) {
	s := SubAfter("asd@86@yes", "@")
	if s != `86@yes` {
		t.Error("s is not `86@yes` , but ", s)
		return
	}
}

func TestSubAfterLast(t *testing.T) {
	s := SubAfterLast("asd@8@yes", "@")
	if s != `yes` {
		t.Error("s is not `yes` , but ", s)
		return
	}
}
