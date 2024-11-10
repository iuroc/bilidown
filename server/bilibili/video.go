package bilibili

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
)

// GetBVInfo 根据 BVID 获取视频信息
func (client *BiliClient) GetVideoInfo(bvid string) (*VideoInfo, error) {
	if client.SESSDATA == "" {
		return nil, errors.New("SESSDATA 不能为空")
	}
	params := map[string]string{"bvid": bvid}
	response, err := client.SimpleGET("https://api.bilibili.com/x/web-interface/wbi/view", params)
	if err != nil {
		return nil, err
	}
	body := BaseResV2{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	bvInfo := VideoInfo{}
	err = json.Unmarshal(body.Data, &bvInfo)
	if err != nil {
		return nil, err
	}
	return &bvInfo, nil
}

// GetSeasonInfo 根据 EPID 或 SSID 获取剧集信息
func (client *BiliClient) GetSeasonInfo(epid int, ssid int) (*SeasonInfo, error) {
	if client.SESSDATA == "" {
		return nil, errors.New("SESSDATA 不能为空")
	}
	params := map[string]string{
		"ep_id":     strconv.Itoa(epid),
		"season_id": strconv.Itoa(ssid),
	}
	response, err := client.SimpleGET("https://api.bilibili.com/pgc/view/web/season", params)
	if err != nil {
		return nil, err
	}
	body := BaseResV3{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	seasonInfo := SeasonInfo{}
	err = json.Unmarshal(body.Result, &seasonInfo)
	if err != nil {
		return nil, err
	}
	return &seasonInfo, nil
}

// GetPlayInfo 根据 BVID 和 CID 获取视频播放信息
func (client *BiliClient) GetPlayInfo(bvid string, cid int) (*PlayInfo, error) {
	if client.SESSDATA == "" {
		return nil, errors.New("SESSDATA 不能为空")
	}
	params := map[string]string{
		"bvid":  bvid,
		"cid":   strconv.Itoa(cid),
		"fnval": "4048",
		"fnver": "0",
		"fourk": "1",
	}
	response, err := client.SimpleGET("https://api.bilibili.com/x/player/playurl", params)
	if err != nil {
		return nil, err
	}
	body := BaseResV2{}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	playInfo := PlayInfo{}
	err = json.Unmarshal(body.Data, &playInfo)
	if err != nil {
		return nil, err
	}
	return &playInfo, nil
}

func (client *BiliClient) GetPopularVideos() ([]VideoInfo, error) {
	if client.SESSDATA == "" {
		return nil, errors.New("SESSDATA 不能为空")
	}
	urls := []string{
		"https://api.bilibili.com/x/web-interface/popular",
		"https://api.bilibili.com/x/web-interface/popular/precious",
		"https://api.bilibili.com/x/web-interface/ranking/v2",
	}
	response, err := client.SimpleGET(urls[rand.Intn(len(urls))], nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body := BaseResV2{}

	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	data := struct {
		List []VideoInfo `json:"list"`
	}{}

	err = json.Unmarshal(body.Data, &data)
	if err != nil {
		return nil, err
	}
	return data.List, nil
}
