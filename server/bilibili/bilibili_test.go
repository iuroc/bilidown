package bilibili_test

import (
	"bilidown/bilibili"
	"encoding/json"
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

	data, err := json.Marshal(videoInfo)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))
}

func TestGetPlayInfo(t *testing.T) {
	for i := 0; i < 20; i++ {
		playInfo, err := bilibili.GetPlayInfo("", "BV1sL4y1G7u7", 865269234, "21adb20f%2C1742230649%2Cef096%2A91CjBTOoXtCMcdxnllvPKj-X5VmTcK_jg0KheA-P1hRjqaS3prbgKKXKrcWYtmX9OsxAsSVmJIaDlVcW1uUEZuY0J5WDQ0S3ByU0VMZTc5S2REUjJpOFpRcDEzWkJreUFuT2RuZ2NVd2VUaVNHLVlVdllJVkNJbWVrRDU5ZXdPVmhlYVlnM3d1NjV3IIEC")
		if err != nil {
			t.Error(err)
		}

		data, err := json.Marshal(playInfo)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(len(string(data)))
	}
}
