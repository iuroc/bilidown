package bilibili_test

import (
	"bilidown/bilibili"
	"fmt"
	"testing"
)

func TestUserHasLogin(t *testing.T) {
	biliUser := bilibili.BiliUser{
		Sessdata: "21adb20f%2C1742230649%2Cef096%2A91CjBTOoXtCMcdxnllvPKj-X5VmTcK_jg0KheA-P1hRjqaS3prbgKKXKrcWYtmX9OsxAsSVmJIaDlVcW1uUEZuY0J5WDQ0S3ByU0VMZTc5S2REUjJpOFpRcDEzWkJreUFuT2RuZ2NVd2VUaVNHLVlVdllJVkNJbWVrRDU5ZXdPVmhlYVlnM3d1NjV3IIEC",
	}

	hasLogin, err := biliUser.HasLogin()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(hasLogin)
}
