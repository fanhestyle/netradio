package medialib

import "testing"

func TestGetChannelInfo(t *testing.T) {
	v := GetChannelInfo()

	if len(v) != 3 {
		t.Fatalf("exptected: %v, actual %v\n", 3, v)
	}

}
