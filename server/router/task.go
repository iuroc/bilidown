package router

import (
	"bilidown/task"
	"bilidown/util"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		util.Res{Success: false, Message: "不支持的请求方法"}.Write(w)
		return
	}
	var body []task.TaskInDB
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	db := util.MustGetDB()
	defer db.Close()
	for _, item := range body {
		if !util.CheckBvidFormat(item.Bvid) {
			util.Res{Success: false, Message: "bvid 格式错误"}.Write(w)
			return
		}
		if item.Cover == "" || item.Title == "" || item.Owner == "" {
			util.Res{Success: false, Message: "参数错误"}.Write(w)
		}

		if !util.IsValidURL(item.Cover) {
			util.Res{Success: false, Message: "封面链接格式错误"}.Write(w)
			return
		}
		if !util.IsValidURL(item.Audio) {
			util.Res{Success: false, Message: "音频链接格式错误"}.Write(w)
			return
		}
		if !util.IsValidURL(item.Video) {
			util.Res{Success: false, Message: "视频链接格式错误"}.Write(w)
			return
		}
		if !util.IsValidFormatCode(item.Format) {
			util.Res{Success: false, Message: "清晰度代码错误"}.Write(w)
			return
		}
		item.Folder, err = util.GetCurrentFolder(db)
		item.Status = "waiting"
		if err != nil {
			util.Res{Success: false, Message: fmt.Sprintf("util.GetCurrentFolder: %v.", err)}.Write(w)
			return
		}
		_task := task.Task{TaskInDB: item}
		_task.Title = util.FilterFileName(_task.Title)
		err = _task.Create(db)
		if err != nil {
			util.Res{Success: false, Message: fmt.Sprintf("_task.Create: %v.", err)}.Write(w)
			return
		}
		go _task.Start()
	}
	util.Res{Success: true, Message: "创建成功"}.Write(w)
}

func GetActiveTask(w http.ResponseWriter, r *http.Request) {
	util.Res{Success: true, Data: task.GlobalTaskList}.Write(w)
}

func GetTaskList(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	db := util.MustGetDB()
	defer db.Close()
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(r.FormValue("pageSize"))
	if err != nil {
		pageSize = 360
	}
	tasks, err := task.GetTaskList(db, page, pageSize)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "获取成功", Data: tasks}.Write(w)
}

// ShowFile 调用 Explorer 查看文件位置
func ShowFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	filePath := r.FormValue("filePath")

	var cmd *exec.Cmd

	// 根据操作系统选择命令
	switch runtime.GOOS {
	case "windows":
		// Windows 使用 explorer
		cmd = exec.Command("explorer", "/select,", filePath)
	case "darwin":
		// macOS 使用 open
		cmd = exec.Command("open", "-R", filePath)
	case "linux":
		// Linux 使用 xdg-open
		cmd = exec.Command("xdg-open", filePath)
	default:
		util.Res{Success: false, Message: "不支持的操作系统"}.Write(w)
		return
	}
	err := cmd.Start()
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "操作成功"}.Write(w)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.FormValue("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		util.Res{Success: false, Message: "参数错误"}.Write(w)
		return
	}
	db := util.MustGetDB()
	defer db.Close()

	_task, err := task.GetTask(db, taskID)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}

	filePath := _task.FilePath()
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}

	err = task.DeleteTask(db, taskID)
	if err != nil {
		util.Res{Success: false, Message: err.Error()}.Write(w)
		return
	}
	util.Res{Success: true, Message: "删除成功"}.Write(w)
}
