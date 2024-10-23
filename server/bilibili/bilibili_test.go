package bilibili_test

import (
	"bilidown/bilibili"
	"fmt"
	"testing"
)

func TestGetPageList(t *testing.T) {
	pageList, err := bilibili.GetPageList("", "BV1fK4y1t7hj")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", pageList)
}

func TestGetVideoInfo(t *testing.T) {
	videoInfo, err := bilibili.GetVideoInfo("", "BV1fK4y1t7hj")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", videoInfo)
}
