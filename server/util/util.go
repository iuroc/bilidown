package util

import (
	"bilidown/common"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	db, err := sql.Open("sqlite", pathStr)
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
func IsValidFormatCode(format common.MediaFormat) bool {
	allowed := []common.MediaFormat{6, 16, 32, 64, 74, 80, 112, 116, 120, 125, 126, 127}
	for _, v := range allowed {
		if v == format {
			return true
		}
	}
	return false
}

// FilterFileName 过滤文件名特殊字符
func FilterFileName(fileName string) string {
	return regexp.MustCompile(`[\\/:*?"<>|\n]`).ReplaceAllString(fileName, "")
}

func GenerateRandomUserAgent() string {
	rand.NewSource(time.Now().UnixNano())
	firstLetter := string(rune(rand.Intn(26) + 'A'))
	lettersLength := rand.Intn(3) + 4
	var sb strings.Builder
	for i := 0; i < lettersLength; i++ {
		sb.WriteByte(byte(rand.Intn(26) + 'a'))
	}
	version := fmt.Sprintf("%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10))
	return fmt.Sprintf("%s%s/%s", firstLetter, sb.String(), version)
}

func GetFFmpegPath() (string, error) {

	if err := exec.Command("ffmpeg", "-version").Run(); err == nil {
		return "ffmpeg", nil
	}

	if err := exec.Command("bin/ffmpeg", "-version").Run(); err == nil {
		return "bin/ffmpeg", nil
	}

	return "", errors.New("ffmpeg not found")
}
