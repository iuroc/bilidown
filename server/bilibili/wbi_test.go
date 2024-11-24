package bilibili_test

import (
	"bilidown/bilibili"
	"bilidown/util"
	"fmt"
	"testing"
	"time"
)

func TestGetWbiKey(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	client := bilibili.BiliClient{SESSDATA: sessdata}
	imgKey, subKey, err := util.GetWbiKey(db)
	if err != nil {
		t.Error(err)
	}
	if imgKey == "" || subKey == "" {
		// 更新数据库中的 WebKey
		if imgKey, subKey, err = client.GetWbiKey(); err != nil {
			t.Error(err)
		}
		if err = util.SaveWbiKey(db, imgKey, subKey); err != nil {
			t.Error(err)
		}
		fmt.Println("更新缓存")
	} else {
		fmt.Println("读取缓存")
	}
	fmt.Println(imgKey, subKey)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}
