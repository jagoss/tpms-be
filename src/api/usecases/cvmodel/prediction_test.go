package cvmodel

import "testing"

func TestMapTo8bitValue(t *testing.T) {
	if mapTo8bitValue(0xffff) != 255 {
		t.Fail()
	}
}
