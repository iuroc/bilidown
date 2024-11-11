package router

import (
	"encoding/base64"
	"net/http"

	"bilidown/bilibili"
	"bilidown/util"

	"github.com/skip2/go-qrcode"
)

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
	db := util.MustGetDB()
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
	db := util.MustGetDB()
	defer db.Close()
	err = bilibili.SaveSessdata(db, sessdata)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "登录成功"}.Write(w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	db := util.MustGetDB()
	defer db.Close()
	err := bilibili.SaveSessdata(db, "")
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "退出成功"}.Write(w)
}
