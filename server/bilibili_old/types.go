package bilibili

import "encoding/json"

type BaseRes struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Ttl     int             `json:"ttl"`
	Data    json.RawMessage `json:"data"`
}

type PageListItem struct {
	Cid       int    `json:"cid"`
	Page      int    `json:"page"`
	From      string `json:"from"`
	Part      string `json:"part"`
	Duration  int    `json:"duration"`
	Dimension `json:"dimension"`
}

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

type VideoInfo struct {
	Bvid    string `json:"bvid"`
	Aid     int    `json:"aid"`
	Pic     string `json:"pic"`
	Title   string `json:"title"`
	Pubdate int    `json:"pubdate"`
	Desc    string `json:"desc"`
	Owner   struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Dimension `json:"dimension"`
	Staff     Staff          `json:"staff"`
	Pages     []PageListItem `json:"pages"`
	Duration  int            `json:"duration"`
	Stat      struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		NowRank  int `json:"now_rank"`
		HisRank  int `json:"his_rank"`
		Like     int `json:"like"`
		Dislike  int `json:"dislike"`
	} `json:"stat"`
}

type Staff []struct {
	Mid   int    `json:"mid"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Face  string `json:"face"`
}

type SeasonVideoInfo struct {
	Actors string `json:"actors"`
	Areas  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"areas"`
	Cover    string `json:"cover"`
	Evaluate string `json:"evaluate"`
}
