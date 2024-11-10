package bilibili_test

import (
	"bilidown/bilibili"
	"bilidown/util"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/skip2/go-qrcode"
	_ "modernc.org/sqlite"
)

const TEST_SESSDATA = "70ca7a2f%2C1745349417%2C51969%2Aa2CjCVKFrX3jG4cXTV2QY3u7NWYDJ9qwHRW5D4ivdRppolKL0a6KL9vmzbrFGcwzN7tSMSVkVuZGtyUVJPNzJGMjEwTTdveGtuZjQ1TU9hSTJSdnA3NHhzTmIxT3dTY3NyYXhRUkxjREVnd0p3NWg2ODh5SVZmZVhlOW9aallOZzN6aVE5M2Y5SldRIIEC"
const TEST_BVID = "BV1KX4y1V7sA"     // 普通合集
const TEST_CID = 305542578           // 普通分集
const TEST_BVID_HDR = "BV1rp4y1e745" // HDR 合集
const TEST_CID_HDR = 244954665       // HDR 分集
const TEST_EPID = 835909             // 剧集
const TEST_SSID = 48744              // 番剧

func TestBiliClient(t *testing.T) {
	client := bilibili.BiliClient{SESSDATA: TEST_SESSDATA}
	check, err := client.CheckLogin()
	if err != nil {
		t.Error(err)
	}
	if check {
		fmt.Println("SESSDATA 有效")
	} else {
		fmt.Println("SESSDATA 无效")
	}
}

func TestNewQRInfo(t *testing.T) {
	client := bilibili.BiliClient{}
	info, err := client.NewQRInfo()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("二维码信息：%+v", info)
}

func TestGetQRStatus(t *testing.T) {
	client := bilibili.BiliClient{}
	qrInfo, err := client.NewQRInfo()
	if err != nil {
		t.Error(err)
	}
	qr, err := qrcode.New(qrInfo.URL, qrcode.Low)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(qr.ToSmallString(false))

	for {
		qrStatus, sessdata, err := client.GetQRStatus(qrInfo.QrcodeKey)
		fmt.Println(qrStatus.Message)
		if err != nil {
			t.Error(err)
		}
		if qrStatus.Code == bilibili.QR_SUCCESS {
			fmt.Println("登录成功")
			fmt.Println(sessdata)
			break
		} else if qrStatus.Code == bilibili.QR_EXPIRES {
			break
		}
		time.Sleep(time.Second)
	}
}

func TestGetBVInfo(t *testing.T) {
	client := bilibili.BiliClient{SESSDATA: TEST_SESSDATA}
	videoInfo, err := client.GetVideoInfo(TEST_BVID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", videoInfo)
}

func TestGetSeasonInfo(t *testing.T) {
	client := bilibili.BiliClient{SESSDATA: TEST_SESSDATA}
	seasonInfo, err := client.GetSeasonInfo(TEST_EPID, TEST_SSID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", seasonInfo)
}

func TestGetPlayInfo(t *testing.T) {
	client := bilibili.BiliClient{SESSDATA: TEST_SESSDATA}
	playInfo, err := client.GetPlayInfo(TEST_BVID_HDR, TEST_CID_HDR)
	if err != nil {
		t.Error(err)
	}
	a, _ := json.Marshal(playInfo.AcceptQuality)
	b, _ := json.Marshal(playInfo.AcceptDescription)
	c, _ := json.Marshal(playInfo.SupportFormats)

	fmt.Println(string(a))
	fmt.Println(string(b))
	fmt.Println(string(c))
}

func TestSaveSessdata(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	err := bilibili.SaveSessdata(db, TEST_SESSDATA)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("SESSDATA 保存成功")
	}
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil {
		t.Error(err)
	}
	if sessdata != TEST_SESSDATA {
		t.Error("SESSDATA 读写不一致")
	} else {
		fmt.Println("SESSDATA 读取成功")
	}
}

func TestGetSessdata(t *testing.T) {
	db := util.MustGetDB("../data.db")
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(sessdata)
}

func TestGetPopularVideos(t *testing.T) {
	client := bilibili.BiliClient{SESSDATA: TEST_SESSDATA}
	videos, err := client.GetPopularVideos()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", videos)
}
