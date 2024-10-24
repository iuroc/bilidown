package router

import (
	"bilidown/bilibili"
	"bilidown/util"
	"net/http"
	"regexp"
)

func API() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/checkLogin", func(w http.ResponseWriter, r *http.Request) {

	})
	router.HandleFunc("/getVideoInfo", GetVideoInfo)
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
	if !regexp.MustCompile("^BV1.+").MatchString(bvid) {
		util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
		return
	}

	videoInfo, err := bilibili.GetVideoInfo("", bvid)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: videoInfo}.Write(w)
}
