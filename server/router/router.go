package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/iuroc/server/bilidown/util"
	"github.com/iuroc/server/bilidown/util/res_error"
)

func API() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/getVideoInfo", GetVideoInfo)
	router.HandleFunc("/getSeasonInfo", GetSeasonInfo)
	router.HandleFunc("/getQRInfo", GetQRInfo)
	router.HandleFunc("/getQRStatus", GetQRStatus)
	router.HandleFunc("/checkLogin", CheckLogin)
	router.HandleFunc("/getPlayInfo", GetPlayInfo)
	router.HandleFunc("/createTask", CreateTask)
	router.HandleFunc("/getActiveTask", GetActiveTask)
	router.HandleFunc("/getTaskList", GetTaskList)
	router.HandleFunc("/showFile", ShowFile)
	router.HandleFunc("/getFields", GetFields)
	router.HandleFunc("/saveFields", SaveFields)
	router.HandleFunc("/logout", Logout)
	router.HandleFunc("/quit", Quit)
	router.HandleFunc("/getPopularVideos", GetPopularVideos)
	router.HandleFunc("/deleteTask", DeleteTask)
	router.HandleFunc("/getRedirectedLocation", GetRedirectedLocation)

	return router
}

func GetRedirectedLocation(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		res_error.ParamError(w)
		return
	}
	url := r.FormValue("url")
	if !util.IsValidURL(url) {
		res_error.URLFormatError(w)
		return
	}
	if location, err := util.GetRedirectedLocation(url); err != nil {
		res_error.NoRedirectedLocation(w)
		return
	} else {
		util.Res{Success: true, Message: "获取成功", Data: location}.Write(w)
		return
	}
}

func Quit(w http.ResponseWriter, r *http.Request) {
	util.Res{Success: true, Message: "退出成功"}.Write(w)
	go func() {
		os.Exit(0)
	}()
}

func GetFields(w http.ResponseWriter, r *http.Request) {
	db := util.MustGetDB()
	defer db.Close()

	fields, err := util.GetFields(db, util.FieldUtil{}.AllowSelect()...)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Data: fields}.Write(w)
}

func SaveFields(w http.ResponseWriter, r *http.Request) {
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
