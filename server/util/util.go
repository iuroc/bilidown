package util

import (
	"bilidown/bilibili"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
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

func IsNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

// IsValidURL 判断字符串是否为合法的URL
func IsValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

// IsValidFormatCode 判断格式码是否合法
func IsValidFormatCode(format bilibili.MediaFormat) bool {
	allowed := []bilibili.MediaFormat{6, 16, 32, 64, 74, 80, 112, 116, 120, 125, 126, 127}
	for _, v := range allowed {
		if v == format {
			return true
		}
	}
	return false
}

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return fmt.Sprintf("%x", randomBytes)[:length]
}

// FilterFileName 过滤文件名特殊字符
func FilterFileName(fileName string) string {
	return regexp.MustCompile(`[\\/:*?"<>|\n]`).ReplaceAllString(fileName, "")
}
