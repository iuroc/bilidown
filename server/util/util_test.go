package util_test

import (
	"bilidown/common"
	"bilidown/util"
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

func TestGenerateRandomUserAgent(t *testing.T) {
	for i := 0; i < 10; i++ {
		ua := util.GenerateRandomUserAgent()
		t.Log(ua)
	}
}
