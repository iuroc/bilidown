package router

import (
	"bilidown/bilibili"
	"bilidown/util"
	"encoding/base64"
	"net/http"
	"regexp"
	"strconv"

	"github.com/skip2/go-qrcode"
)

func API() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/checkLogin", func(w http.ResponseWriter, r *http.Request) {

	})
	router.HandleFunc("/getVideoInfo", GetVideoInfo)
	router.HandleFunc("/getSeasonInfo", GetSeasonInfo)
	router.HandleFunc("/getQRInfo", GetQRInfo)
	return router
}

func GetVideoInfo(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	bvid := r.FormValue("bvid")
	if bvid == "" {
		util.Res{Success: false, Message: "bvid 不能为空"}.Write(w)
		return
	}
	if !regexp.MustCompile("^BV1[a-zA-Z0-9]+").MatchString(bvid) {
		util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()

	sessdata, err := bilibili.GetSessdata(db)
	if err != nil || sessdata == "" {
		util.Res{Success: false, Message: "未登录"}.Write(w)
		return
	}
	client := bilibili.BiliClient{SESSDATA: sessdata}
	videoInfo, err := client.GetVideoInfo(bvid)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: videoInfo}.Write(w)
}

func GetSeasonInfo(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	epid, err := strconv.Atoi(r.FormValue("epid"))
	if err != nil {
		util.Res{Success: false, Message: "epid 格式错误"}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil || sessdata == "" {
		util.Res{Success: false, Message: "未登录"}.Write(w)
		return
	}

	client := bilibili.BiliClient{SESSDATA: sessdata}
	seasonInfo, err := client.GetSeasonInfo(epid)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: seasonInfo}.Write(w)
}

func GetQRInfo(w http.ResponseWriter, r *http.Request) {
	client := bilibili.BiliClient{}
	qrInfo, err := client.NewQRInfo()
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	imageData, err := qrcode.Encode(qrInfo.URL, qrcode.Medium, 256)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	w.Header().Set("Cache-Control", "no-store")
	util.Res{
		Success: true,
		Message: "获取成功",
		Data: struct {
			Key   string `json:"key"`
			Image string `json:"image"`
		}{
			Key:   qrInfo.QrcodeKey,
			Image: "data:image/png;base64," + base64Str,
		}}.Write(w)
}
