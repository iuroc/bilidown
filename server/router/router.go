package router

import (
	"bilidown/util"
	"github.com/sqweek/dialog"
	"net/http"
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

// FolderPicker 调用资源管理器选择文件夹
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
