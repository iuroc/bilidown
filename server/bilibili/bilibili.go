package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// GetVideoInfo 获取视频详细信息，包括合集内视频列表。
func GetVideoInfo(avid string, bvid string) (videoInfo *VideoInfo, err error) {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("aid", avid)
	_url := "https://api.bilibili.com/x/web-interface/wbi/view?" + params.Encode()
	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "github@iuroc")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := BaseRes{}
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return videoInfo, errors.New(body.Message)
	}
	videoInfo = &VideoInfo{}
	err = json.Unmarshal(body.Data, &videoInfo)
	if err != nil {
		return videoInfo, err
	}
	return videoInfo, nil
}

type Media struct {
	ID        int      `json:"id"`
	BaseURL   string   `json:"baseUrl"`
	BackupURL []string `json:"backupUrl"`
	Bandwidth int      `json:"bandwidth"`
	MimeType  string   `json:"mimeType"`
	Codecs    string   `json:"codecs"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	FrameRate string   `json:"frameRate"`
	Codecid   int      `json:"codecid"`
}

type PlayInfo struct {
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	SupportFormats    []struct {
		Quality        int      `json:"quality"`
		Format         string   `json:"format"`
		NewDescription string   `json:"new_description"`
		Codecs         []string `json:"codecs"`
	} `json:"support_formats"`
	Dash struct {
		Duration int     `json:"duration"`
		Video    []Media `json:"video"`
		Audio    []Media `json:"audio"`
	}
}

// GetPlayInfo 获取音视频播放地址信息。
func GetPlayInfo(avid string, bvid string, cid int, sessdata string) (playInfo *PlayInfo, err error) {
	params := url.Values{}
	params.Set("avid", avid)
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(cid))
	params.Set("fnval", "4048")
	params.Set("fnver", "0")
	params.Set("fourk", "1")
	_url := "https://api.bilibili.com/x/player/wbi/playurl?" + params.Encode()
	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Cookie", "SESSDATA="+sessdata)
	request.Header.Set("User-Agent", "github@iuroc")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := BaseRes{}
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	fmt.Println(string(body.Data))

	playInfo = &PlayInfo{}
	err = json.Unmarshal(body.Data, &playInfo)
	if err != nil {
		return nil, err
	}
	return playInfo, nil
}
