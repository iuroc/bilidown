package router

import (
	"bilidown/bilibili"
	"bilidown/util"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/skip2/go-qrcode"
	"github.com/sqweek/dialog"
)

func API() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/getVideoInfo", GetVideoInfo)
	router.HandleFunc("/getSeasonInfo", GetSeasonInfo)
	router.HandleFunc("/getQRInfo", GetQRInfo)
	router.HandleFunc("/getQRStatus", GetQRStatus)
	router.HandleFunc("/checkLogin", CheckLogin)
	router.HandleFunc("/getPlayInfo", GetPlayInfo)
	router.HandleFunc("/folderPicker", FolderPicker)
	return router
}

func FolderPicker(w http.ResponseWriter, r *http.Request) {
	folderPath, err := dialog.Directory().Title("Select a directory").Browse()
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ("save_folder", ?)`, folderPath)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "选择成功", Data: folderPath}.Write(w)
}

func GetVideoInfo(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	bvid := r.FormValue("bvid")
	if !util.CheckBVID(bvid) {
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
	var epid int
	epid, err := strconv.Atoi(r.FormValue("epid"))
	if r.FormValue("epid") != "" && err != nil {
		util.Res{Success: false, Message: "epid 格式错误"}.Write(w)
		return
	}
	var ssid int
	if epid == 0 {
		ssid, err = strconv.Atoi(r.FormValue("ssid"))
		if r.FormValue("ssid") != "" && err != nil {
			util.Res{Success: false, Message: "ssid 格式错误"}.Write(w)
			return
		}
	}
	db := util.GetDB()
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil || sessdata == "" {
		util.Res{Success: false, Message: "未登录"}.Write(w)
		return
	}

	client := bilibili.BiliClient{SESSDATA: sessdata}
	seasonInfo, err := client.GetSeasonInfo(epid, ssid)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: seasonInfo}.Write(w)
}

func GetQRInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
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

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	db := util.GetDB()
	defer db.Close()
	sessdata, err := bilibili.GetSessdata(db)
	if err != nil || sessdata == "" {
		util.Res{Success: false, Message: "未登录"}.Write(w)
		return
	}
	client := bilibili.BiliClient{SESSDATA: sessdata}
	check, err := client.CheckLogin()
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	if check {
		util.Res{Success: true, Message: "登录成功"}.Write(w)
	} else {
		util.Res{Success: false, Message: "登录失败"}.Write(w)
	}
}

// GetPlayInfo 获取视频播放信息
func GetPlayInfo(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}

	bvid := r.FormValue("bvid")
	if !util.CheckBVID(bvid) {
		util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
		return
	}
	cid, err := strconv.Atoi(r.FormValue("cid"))
	if err != nil {
		util.Res{Success: false, Message: "cid 格式错误"}.Write(w)
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
	playInfo, err := client.GetPlayInfo(bvid, cid)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: playInfo}.Write(w)
}

// GetQRStatus 获取二维码状态
func GetQRStatus(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}

	key := r.FormValue("key")
	if key == "" {
		util.Res{Success: false, Message: "key 不能为空"}.Write(w)
		return
	}
	client := bilibili.BiliClient{}
	qrStatus, sessdata, err := client.GetQRStatus(key)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	if qrStatus.Code != bilibili.QR_SUCCESS {
		util.Res{Success: false, Message: qrStatus.Message}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()
	err = bilibili.SaveSessdata(db, sessdata)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "登录成功"}.Write(w)
}
