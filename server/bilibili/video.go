package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

// GetSeasonsArchivesList 获取合集中的第一个视频的 BVID
func (client *BiliClient) GetSeasonsArchivesListFirstBvid(mid int, seasonId int) (string, error) {
	url := "https://api.bilibili.com/x/polymer/web-space/seasons_archives_list"
	params := map[string]string{
		"mid":       strconv.Itoa(mid),
		"season_id": strconv.Itoa(seasonId),
		"page_num":  "1",
		"page_size": "1",
	}

	response, err := client.SimpleGET(url, params)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body := BaseResV2{}
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Code != 0 {
		return "", errors.New(body.Message)
	}
	var data struct {
		Archives []struct {
			Bvid string `json:"bvid"`
		} `json:"archives"`
	}
	if err = json.Unmarshal(body.Data, &data); err != nil {
		return "", nil
	}
	if len(data.Archives) == 0 {
		return "", errors.New("视频列表为空")
	}
	return data.Archives[0].Bvid, nil
}

func (client *BiliClient) GetFavlist(mediaId int) (*FavList, error) {
	if client.SESSDATA == "" {
		return nil, errors.New("SESSDATA 不能为空")
	}
	page := 0
	retry := 0
	allFavList := FavList{}
	for {
		favList, hasMore, err := client.GetFavlistByPage(mediaId, page, 40)
		if err != nil {
			fmt.Println(err.Error())
			if retry == 5 || strings.HasPrefix(err.Error(), "body.Code not 0") {
				return nil, err
			}
			retry++
			continue
		}
		allFavList = append(allFavList, *favList...)
		if !hasMore {
			break
		}
		page++
	}
	return &allFavList, nil
}

func (client *BiliClient) GetFavlistByPage(mediaId int, page int, pageSize int) (favlist *FavList, hasMore bool, err error) {
	if client.SESSDATA == "" {
		return nil, false, errors.New("SESSDATA 不能为空")
	}
	response, err := client.SimpleGET("https://api.bilibili.com/x/v3/fav/resource/list", map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"pn":       strconv.Itoa(page + 1),
		"ps":       strconv.Itoa(pageSize),
		"order":    "mtime",
		"type":     "0",
		"tid":      "0",
		"platform": "web",
	})
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()
	body := BaseResV2{}
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, false, err
	}
	if body.Code != 0 {
		return nil, false, fmt.Errorf("body.Code not 0, %s", body.Message)
	}
	data := struct {
		Medias  FavList `json:"medias"`
		HasMore bool    `json:"has_more"`
	}{}
	if err = json.Unmarshal(body.Data, &data); err != nil {
		return nil, false, err
	}
	return &data.Medias, data.HasMore, nil
}
