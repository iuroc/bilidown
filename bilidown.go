package bilidown

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

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
	if cookie.Name == "SESSDATA" && expiresToTime(cookie.Expires).After(time.Now()) {
		return cookie.Value, nil
	} else {
		return "", errors.New("无可用 Cookie 或 Cookie 过期")
	}
}

// ExpiresToTime 将 network.Cookie.Expires 转换为 Time
func expiresToTime(expires float64) time.Time {
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
	Desc    string `json:"desc"`
	Pubdate int    `json:"pubdate"`
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

func Download(video VideoItem, audios []AudioItem, dirPath string) {

}
