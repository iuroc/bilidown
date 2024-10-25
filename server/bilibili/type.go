package bilibili

import "encoding/json"

// BaseRes 来自 Bilibili 的接口响应，Message 字段为 msg
type BaseRes struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

// BaseRes2 来自 Bilibili 的接口响应，Message 字段为 message 而不是 msg
type BaseResV2 struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// BaseResV3 来自 Bilibili 的接口响应，Message 字段为 message 而不是 msg，Data 字段为 result
type BaseResV3 struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

func (b *BaseRes) Success() bool {
	return b.Code == 0
}

func (b *BaseResV2) Success() bool {
	return b.Code == 0
}

func (b *BaseResV3) Success() bool {
	return b.Code == 0
}

type QRInfo struct {
	URL       string `json:"url"`
	QrcodeKey string `json:"qrcode_key"`
}

type QRStatus struct {
	URL          string `json:"string"`
	RefreshToken string `json:"refresh_token"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}

const (
	QR_NO_SCAN  = 86101 // 未扫码
	QR_NO_CLICK = 86090 // 已扫码未确认
	QR_EXPIRES  = 86038 // 已过期
	QR_SUCCESS  = 0     // 已确认登录
)

// Dimension 视频分辨率
type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

// StaffItem 视频人员
type StaffItem struct {
	Mid   int    `json:"mid"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Face  string `json:"face"`
}

// Page 分集信息
type Page struct {
	// 配合 Bvid 用于获取播放地址
	Cid int `json:"cid"`
	// 当前分集在合集中的序号，从 1 开始
	Page int `json:"page"`
	// 分集标题
	Part string `json:"part"`
	// 分集时长
	Duration int `json:"duration"`
	// 分集分辨率
	Dimension `json:"dimension"`
}

// 通过 BVID 获取的视频信息
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
	Staff     []StaffItem `json:"staff"`
	Pages     []Page      `json:"pages"`
	Duration  int         `json:"duration"`
	Stat      struct {
		Aid      int `json:"aid"`
		View     int `json:"view"`     // 播放数量
		Danmaku  int `json:"danmaku"`  // 弹幕数量
		Reply    int `json:"reply"`    // 评论数量
		Favorite int `json:"favorite"` // 收藏数量
		Coin     int `json:"coin"`     // 投币数量
		Share    int `json:"share"`    // 转发数量
		NowRank  int `json:"now_rank"` // 当前排名
		HisRank  int `json:"his_rank"` // 历史最高排名
		Like     int `json:"like"`     // 点赞数量
		Dislike  int `json:"dislike"`  // 不喜欢
	} `json:"stat"`
}

// Episode 剧集分集信息
type Episode struct {
	Aid       int                `json:"aid"`
	Bvid      string             `json:"bvid"`
	Cid       int                `json:"cid"`
	Cover     string             `json:"cover"` // 封面
	Dimension `json:"dimension"` // 分辨率
	Duration  int                `json:"duration"` // 时长
	EPID      int                `json:"ep_id"`
	LongTitle string             `json:"long_title"` // 分集完整标题，比如【法外狂徒张三现身！】
	PubTime   int                `json:"pub_time"`   // 发布时间
	Title     string             `json:"title"`      // 分集简略标题，比如【1】
}

// 通过 EPID 获取的视频信息
type SeasonInfo struct {
	Actors string `json:"actors"` // 演员名单
	Areas  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"areas"` // 地区列表
	Cover    string `json:"cover"`    // 封面
	Evaluate string `json:"evaluate"` // 简介
	Publish  struct {
		IsFinish int    `json:"is_finish"` // 是否完结
		PubTime  string `json:"pub_time"`  // 发布时间
	} `json:"publish"`
	SeasonID    int    `json:"season_id"`    // 剧集编号
	SeasonTitle string `json:"season_title"` // 剧集标题
	Stat        struct {
		Coins     int `json:"coins"`     // 投币数量
		Danmakus  int `json:"danmakus"`  // 弹幕数量
		Favorite  int `json:"favorite"`  // 收藏数量
		Favorites int `json:"favorites"` // 追剧数量
		Likes     int `json:"likes"`     // 点赞数量
		Reply     int `json:"reply"`     // 评论数量
		Share     int `json:"share"`     // 分享数量
		Views     int `json:"views"`     // 播放数量
	} `json:"stat"`
	Styles []string `json:"styles"` // 剧集内容类型，例如 [ "短剧", "奇幻", "搞笑" ]
	Title  string   `json:"title"`  // 剧集标题
	Total  int      `json:"total"`  // 总集数

	Episodes []Episode `json:"episodes"` // 分集信息列表

	NewEp struct {
		Desc  string `json:"desc"`   // 更新状态文本
		IsNew int    `json:"is_new"` // 是否是连载，0 为完结，1 为连载
	} `json:"new_ep"`
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
	} `json:"dash"`
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
