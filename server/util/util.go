package util

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// 统一的 JSON 响应结构
type Res struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 发送响应
func (r Res) Write(w http.ResponseWriter) {
	bs, err := json.Marshal(r)
	if err != nil {
		w.Write([]byte(`{"success":false,"message":"系统错误","data":null}`))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func GetDB(path ...string) *sql.DB {
	pathStr := ""
	if len(path) == 0 {
		pathStr = "./data.db"
	} else {
		pathStr = path[0]
	}
	db, err := sql.Open("sqlite3", pathStr)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CheckBVID(bvid string) bool {
	return regexp.MustCompile("^BV1[a-zA-Z0-9]+").MatchString(bvid)
}

// GetDefaultDownloadFolder 获取默认下载路径
func GetDefaultDownloadFolder() (string, error) {
	return filepath.Abs("./download")
}

// SaveDownloadFolder 保存下载路径，不存在则自动创建
func SaveDownloadFolder(db *sql.DB, downloadFolder string) error {
	_, err := os.Stat(downloadFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloadFolder, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return err
	}
	_, err = db.Exec(`INSERT OR REPLACE INTO "field" ("name", "value") VALUES ("download_folder", ?)`, downloadFolder)
	return err
}
