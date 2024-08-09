package bilidown

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Login 调用浏览器登录并返回 SESSDATA
func Login() (*network.Cookie, error) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.Flag("headless", false),
	)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://passport.bilibili.com/login"),
	)
	if err != nil {
		return nil, err
	}
	var loginCookie *network.Cookie
	for {
		time.Sleep(time.Second)
		err = chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				cookies, err := network.GetCookies().Do(ctx)
				for _, cookie := range cookies {
					if cookie.Name == "SESSDATA" {
						loginCookie = cookie
						return nil
					}
				}
				return err
			}),
		)
		if err != nil {
			return nil, err
		}
		if loginCookie != nil {
			break
		}
	}
	return loginCookie, nil
}

// SaveCookie 将 Cookie 以 JSON 格式保存到 cookie 文件中
func SaveCookie(cookie *network.Cookie, cookieSavePath string) {
	result, err := json.Marshal(cookie)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(cookieSavePath, result, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

// GetCookieValue 获取文件中保存的可用 Cookie
func GetCookieValue(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return "", err
	} else if err != nil {
		log.Fatalln(err)
	}
	var cookie network.Cookie
	err = json.Unmarshal(data, &cookie)
	if err != nil {
		return "", errors.New("cookie 文件内容格式错误")
	}
	if cookie.Name == "SESSDATA" && ExpiresToTime(cookie.Expires).After(time.Now()) {
		return cookie.Value, nil
	} else {
		return "", errors.New("无可用 Cookie 或 Cookie 过期")
	}
}

// ExpiresToTime 将 network.Cookie.Expires 转换为 Time
func ExpiresToTime(expires float64) time.Time {
	seconds := int64(expires)
	nanos := int64((expires - float64(seconds)) * 1e9)
	return time.Unix(seconds, nanos)
}

// CheckVideoURLOrID 校验视频链接或视频 ID 格式
func CheckVideoURLOrID(urlOrId string) (videoId string, err error) {
	match := regexp.MustCompile(`^(?:(?:https?://)?www.bilibili.com/video/)?(BV1[a-zA-Z0-9]+)`)
	result := match.FindStringSubmatch(urlOrId)
	if len(result) == 0 {
		return "", errors.New("视频链接或视频 ID 格式错误")
	} else {
		return result[1], nil
	}
}

// MakeVideoURL 根据视频 ID 构建视频链接
func MakeVideoURL(videoId string) string {
	return "https://www.bilibili.com/video/" + videoId + "/"
}

// ParseVideo 解析视频下载地址
func ParseVideo(videoURL string, cookieValue string) (*ParseResult, error) {
	request, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		return nil, err
	}
	request.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: cookieValue,
	})
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	html := string(bs)
	match := regexp.MustCompile(`window.__playinfo__=(.*?)</script>`)
	result := match.FindStringSubmatch(html)
	if len(result) == 0 {
		return nil, errors.New("响应体内容格式异常")
	}
	playInfoStr := result[1]
	var playInfo PlayInfo
	err = json.Unmarshal([]byte(playInfoStr), &playInfo)
	if err != nil {
		return nil, err
	}
	initialStateMatch := regexp.MustCompile(`window.__INITIAL_STATE__=(.*?});`)
	initResult := initialStateMatch.FindStringSubmatch(html)
	if len(initResult) == 0 {
		return nil, errors.New("响应体内容格式异常")
	}
	initialStateStr := initResult[1]
	var initialState InitialState
	err = json.Unmarshal([]byte(initialStateStr), &initialState)
	if err != nil {
		os.WriteFile("xxx.json", []byte(initialStateStr), 0600)
		return nil, err
	}
	return &ParseResult{
		PlayInfoData: playInfo.Data,
		VideoData:    initialState.Data,
	}, nil
}

type ParseResult struct {
	PlayInfoData
	VideoData
}

type PlayInfo struct {
	Data PlayInfoData `json:"data"`
}

type PlayInfoData struct {
	SupportFormats []FormatItem `json:"support_formats"`
	Dash           struct {
		Audio []AudioItem `json:"audio"`
		Video []VideoItem `json:"video"`
	} `json:"dash"`
}

type InitialState struct {
	Data VideoData `json:"videoData"`
}

type VideoData struct {
	// 视频标题
	Title string `json:"title"`
	// 创作团队
	Staff []StaffItem `json:"staff"`
	// 视频统计信息
	Stat Stat `json:"stat"`
	// 视频描述
	Desc string `json:"desc"`
	// 发布时间
	Pubdate int `json:"pubdate"`
	Owner   struct {
		Name string `json:"name"`
	} `json:"owner"`
}

func (video VideoData) PubdateString() string {
	return time.Unix(int64(video.Pubdate), 0).Format(time.DateOnly)
}

type Stat struct {
	// 投币数量
	Coin int `json:"coin"`
	// 弹幕条数
	Danmaku int `json:"danmaku"`
	// 收藏数量
	Favorite int `json:"favorite"`
	// 点赞数量
	Like int `json:"like"`
	// 分享数量
	Share int `json:"share"`
	// 播放量
	View int `json:"view"`
	// 评论数量
	Reply int `json:"reply"`
}

type StaffItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type FormatItem struct {
	// 视频分辨率 ID
	Quality int `json:"quality"`
	// 分辨率描述
	Description string `json:"new_description"`
	// 编解码器
	Codecs []string `json:"codecs"`
}

type AudioItem struct {
	// 下载地址
	BaseUrl string `json:"baseUrl"`
	// 备用下载地址
	BackupUrl []string `json:"backupUrl"`
	// 编解码器
	Codecs string `json:"codecs"`
	// 比特率
	Bandwidth int `json:"bandwidth"`
}

type VideoItem struct {
	// 视频分辨率 ID
	Id int `json:"id"`
	// 下载地址
	BaseUrl string `json:"baseUrl"`
	// 备用下载地址
	BackupUrl []string `json:"backupUrl"`
	// 编解码器
	Codecs string `json:"codecs"`
	// 比特率
	Bandwidth int `json:"bandwidth"`
}

// Description 视频分辨率描述
func (video VideoItem) Description(formats []FormatItem) string {
	for _, format := range formats {
		if format.Quality == video.Id {
			return format.Description
		}
	}
	return ""
}

// Download 下载并合并指定的视频和音频，保存到输出目录
func Download(parseResult *ParseResult, index int, downDirPath string, tempDirPath string) (outputPath string, err error) {
	var bestAudio AudioItem
	for _, audio := range parseResult.Dash.Audio {
		if audio.Bandwidth > bestAudio.Bandwidth {
			bestAudio = audio
		}
	}
	video := parseResult.Dash.Video[index]
	ClearDir(tempDirPath)
	outputFileName := fmt.Sprintf("%s-%s.mp4", parseResult.Title, parseResult.Owner.Name)
	outputFileName = sanitizeFileName(outputFileName)
	outputPath = filepath.Join(downDirPath, outputFileName)
	fmt.Println("🚩 正在下载视频...")
	tempVideoPath := filepath.Join(tempDirPath, "video")
	err = DownloadFile(video.BaseUrl, tempVideoPath)
	if err != nil {
		return "", err
	}
	fmt.Print("\n\n")
	fmt.Println("🚩 正在下载音频...")
	tempAudioPath := filepath.Join(tempDirPath, "audio")
	err = DownloadFile(bestAudio.BaseUrl, tempAudioPath)
	if err != nil {
		return "", err
	}
	ffmpegExecPath := "ffmpeg"
	if FileExists("./ffmpeg.exe") || FileExists("./ffmpeg") {
		ffmpegExecPath = "./" + ffmpegExecPath
	}
	cmd := exec.Command(ffmpegExecPath, "-i", tempVideoPath, "-i", tempAudioPath, "-vcodec", "copy", "-acodec", "copy", outputPath, "-y")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("合并音视频失败: %v\n%s", err, output)
	}
	ClearDir(tempDirPath)
	return outputPath, nil
}

func DownloadFile(url string, path string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Referer", "https://www.bilibili.com")
	request.Header.Set("User-Agent", "iuroc")
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	totalSize := response.ContentLength
	var downloaded int64
	buffer := make([]byte, 1024)
	for {
		n, err := response.Body.Read(buffer)
		if n > 0 {
			_, writeErr := outFile.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			percent := float64(downloaded) / float64(totalSize) * 100
			fmt.Printf("\r下载进度: %.2f%% (%d/%d bytes)", percent, downloaded, totalSize)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// InitDir 初始化文件夹，如果不存在则自动创建
func InitDir(path string) {
	// 检查文件夹是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件夹不存在，则创建它
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("❗️ 创建文件夹失败: %v\n", err)
			return
		}
		fmt.Printf("✅ 文件夹 '%s' 已成功创建\n", path)
	}
}

// ClearDir 清空文件夹
func ClearDir(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	names, err := file.Readdirnames(-1)
	if err != nil {
		log.Fatalln(err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// sanitizeFileName 处理字符串，去除或替换不适合作为文件名的字符
func sanitizeFileName(filename string) string {
	// 定义不允许出现在文件名中的字符集
	invalidChars := `\/:*?"<>|`

	// 使用 strings.Map 替换敏感字符
	sanitized := strings.Map(func(r rune) rune {
		if strings.ContainsRune(invalidChars, r) || unicode.IsControl(r) {
			return '_' // 用下划线替换敏感字符
		}
		return r
	}, filename)

	// 去除首尾的空白字符和特定特殊字符
	sanitized = strings.TrimSpace(sanitized)
	sanitized = strings.Trim(sanitized, ".")

	return sanitized
}
