package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/iuroc/bilidown/server/common"
	"github.com/iuroc/bilidown/server/util"
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

func TestGetRedirectedLocation(t *testing.T) {
	os.Setenv("https_proxy", "http://192.168.1.5:9000")
	if location, err := util.GetRedirectedLocation("https://b23.tv/Ga6sbzT"); err != nil {
		t.Error(err)
	} else {
		fmt.Println(location)
	}
}
