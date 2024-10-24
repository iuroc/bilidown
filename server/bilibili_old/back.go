package bilibili

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GetPageList(aid string, bvid string) (pageList []PageListItem, err error) {
	url := "https://api.bilibili.com/x/player/pagelist?aid=" + aid + "&bvid=" + bvid
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := BaseRes{}
	pageList = []PageListItem{}

	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, errors.New(body.Message)
	}
	err = json.Unmarshal(body.Data, &pageList)
	if err != nil {
		return nil, err
	}
	return pageList, nil
}
