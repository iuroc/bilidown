package bilibili

import (
	"encoding/json"
	"io"
	"net/http"
)

// BiliUser 模拟访问者，携带 SESSDATA 进行接口调用。
type BiliUser struct {
	Sessdata string
}

// HasLogin 检查当前的 SESSDATA 是否已经登录。[已检查]
func (user *BiliUser) HasLogin() (bool, error) {
	request, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/myinfo", nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Cookie", "SESSDATA="+user.Sessdata)
	request.Header.Set("User-Agent", "github@iuroc")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	body := BaseRes{}
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		return false, err
	}
	return body.Code == 0, nil
}
