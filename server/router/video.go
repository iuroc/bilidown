package router

import (
	"bilidown/bilibili"
	"bilidown/util"
	"net/http"
	"strconv"
)

// GetVideoInfo 通过 BV 号获取视频信息
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

// GetSeasonInfo 通过 EP 号或 SS 号获取视频信息
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

// GetPlayInfo 通过 BVID 和 CID 获取视频播放信息
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
	// playInfo.Dash = nil
	util.Res{Success: true, Message: "获取成功", Data: playInfo}.Write(w)
}
