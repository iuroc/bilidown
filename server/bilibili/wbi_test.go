package bilibili

import (
	"bilidown/util"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetWbiKey(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	client := BiliClient{SESSDATA: sessdata}
	mixinKey, err := client.GetMixinKey(db)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(mixinKey)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestURLEncode(t *testing.T) {
	str := url.Values{
		"foo": {"one one four"},
		"bar": {"五一四"},
		"baz": {"1919810"},
	}.Encode()
	str = strings.ReplaceAll(str, "+", "%20")
	fmt.Println(str)
}

func TestWbiSign(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	client := BiliClient{SESSDATA: sessdata}
	mixinKey, err := client.GetMixinKey(db)
	if err != nil {
		t.Error(err)
	}
	newParams := WbiSign(map[string]string{
		"foo": "114",
		"bar": "514",
		"zab": "1919810",
	}, mixinKey)

	fmt.Printf("%+v", newParams)
}
