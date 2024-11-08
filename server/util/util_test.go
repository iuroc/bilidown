package util_test

import (
	"bilidown/common"
	"testing"
)

func TestRandomString(t *testing.T) {
	for i := 4; i < 10; i++ {
		for j := 0; j < 3; j++ {
			str := common.RandomString(i)
			t.Log(str)
		}
		t.Log("\n")
	}
}