package router

import (
	"bilidown/task"
	"bilidown/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	router.HandleFunc("/createTask", CreateTask)
	router.HandleFunc("/getActiveTask", GetActiveTask)
	router.HandleFunc("/getTaskList", GetTaskList)
	router.HandleFunc("/showFile", ShowFile)
	router.HandleFunc("/getFields", GetFields)
	router.HandleFunc("/saveFields", SaveFields)
	router.HandleFunc("/logout", Logout)
	router.HandleFunc("/quit", Quit)
	return router
}

func Quit(w http.ResponseWriter, r *http.Request) {
	util.Res{Success: true, Message: "退出成功"}.Write(w)
	go func() {
		os.Exit(0)
	}()
}

// FolderPicker 调用资源管理器选择文件夹
func FolderPicker(w http.ResponseWriter, r *http.Request) {
	folderPath, err := dialog.Directory().Title("您希望下载视频到哪个文件夹？").Browse()
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	db := util.GetDB()
	defer db.Close()
	err = task.SaveDownloadFolder(db, folderPath)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "选择成功", Data: folderPath}.Write(w)
}

func GetFields(w http.ResponseWriter, r *http.Request) {
	db := util.GetDB()
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

	db := util.GetDB()
	defer db.Close()

	fu := util.FieldUtil{}

	for _, d := range body {
		if !fu.IsAllowUpdate(d[0]) {
			util.Res{Success: false, Message: fmt.Sprintf("字段 %s 不允许修改", d[0])}.Write(w)
			return
		}

		if d[0] == "download_folder" {
			if _, err := os.Stat(d[1]); os.IsNotExist(err) {
				util.Res{Success: false, Message: fmt.Sprintf("文件夹 %s 不存在", d[1])}.Write(w)
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
