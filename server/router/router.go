package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bilidown/util"
	"bilidown/util/res_error"
)

func API() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/getVideoInfo", getVideoInfo)
	router.HandleFunc("/getSeasonInfo", getSeasonInfo)
	router.HandleFunc("/getQRInfo", getQRInfo)
	router.HandleFunc("/getQRStatus", getQRStatus)
	router.HandleFunc("/checkLogin", checkLogin)
	router.HandleFunc("/getPlayInfo", getPlayInfo)
	router.HandleFunc("/createTask", createTask)
	router.HandleFunc("/getActiveTask", getActiveTask)
	router.HandleFunc("/getTaskList", getTaskList)
	router.HandleFunc("/showFile", showFile)
	router.HandleFunc("/getFields", getFields)
	router.HandleFunc("/saveFields", saveFields)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/quit", quit)
	router.HandleFunc("/getPopularVideos", getPopularVideos)
	router.HandleFunc("/deleteTask", deleteTask)
	router.HandleFunc("/getRedirectedLocation", getRedirectedLocation)
	router.HandleFunc("/downloadVideo", downloadVideo)
	router.HandleFunc("/getSeasonsArchivesListFirstBvid", getSeasonsArchivesListFirstBvid)
	router.HandleFunc("/getFavList", getFavList)
	return router
}

func getRedirectedLocation(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		res_error.Send(w, res_error.ParamError)
		return
	}
	url := r.FormValue("url")
	if !util.IsValidURL(url) {
		res_error.Send(w, res_error.URLFormatError)
		return
	}
	if location, err := util.GetRedirectedLocation(url); err != nil {
		res_error.Send(w, res_error.NoLocationError)
		return
	} else {
		util.Res{Success: true, Message: "获取成功", Data: location}.Write(w)
		return
	}
}

func quit(w http.ResponseWriter, r *http.Request) {
	util.Res{Success: true, Message: "退出成功"}.Write(w)
	go func() {
		os.Exit(0)
	}()
}

func getFields(w http.ResponseWriter, r *http.Request) {
	db := util.MustGetDB()
	defer db.Close()

	fields, err := util.GetFields(db, util.FieldUtil{}.AllowSelect()...)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Data: fields}.Write(w)
}

func saveFields(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Res{Success: false, Message: "不支持的请求方法"}.Write(w)
		return
	}
	defer r.Body.Close()
	var body [][2]string

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}

	db := util.MustGetDB()
	defer db.Close()

	fu := util.FieldUtil{}

	for _, d := range body {
		if !fu.IsAllowUpdate(d[0]) {
			util.Res{Success: false, Message: fmt.Sprintf("字段 %s 不允许修改", d[0])}.Write(w)
			return
		}

		if d[0] == "download_folder" {
			if _, err := os.Stat(d[1]); os.IsNotExist(err) {
				if err := os.MkdirAll(d[1], os.ModePerm); err != nil {
					util.Res{Success: false, Message: fmt.Sprintf("目录创建失败：%s", d[1])}.Write(w)
					return
				}
			} else if err != nil {
				util.Res{Success: false, Message: fmt.Sprintf("路径设置失败：%v", err)}.Write(w)
				return
			}
		}
	}

	err = util.SaveFields(db, body)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "保存成功"}.Write(w)
}
