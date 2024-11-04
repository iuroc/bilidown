package router

import (
	"bilidown/task"
	"bilidown/util"
	"net/http"

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
	return router
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
	fields, err := util.GetFields(db, "folder")
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Data: fields}.Write(w)
}
