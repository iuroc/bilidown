package util

import (
	"errors"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/iuroc/server/bilidown/common"
)

func CheckBvidFormat(bvid string) bool {
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

// FilterFileName 过滤字符串中的特殊字符，使其允许作为文件名。
func FilterFileName(fileName string) string {
	return regexp.MustCompile(`[\\/:*?"<>|\n]`).ReplaceAllString(fileName, "")
}

// GetFFmpegPath 获取可用的 FFmpeg 执行路径。
func GetFFmpegPath() (string, error) {
	if err := exec.Command("ffmpeg", "-version").Run(); err == nil {
		return "ffmpeg", nil
	}
	if err := exec.Command("bin/ffmpeg", "-version").Run(); err == nil {
		return "bin/ffmpeg", nil
	}
	return "", errors.New("ffmpeg not found")
}

// GetRedirectedLocation 获取响应头中的 Location，不会自动跟随重定向。
func GetRedirectedLocation(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if locationURL, err := response.Location(); err != nil {
		return "", err
	} else {
		return locationURL.String(), nil
	}
}
