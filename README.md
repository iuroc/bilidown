# bilidown

Bilibili 视频下载工具。

## 功能概述

1. 首先检查程序路径下的 Cookie 保存，检查过期状态，如果过期，则提示是否登录，可以选择 n 则游客访问，如果选择 y 则自动弹出浏览器要求登录，登录成功后自动关闭浏览器，并保存 Cookie 到 cookie 文件
2. 提示输入视频地址，提示选择清晰度
3. 调用 ffmpeg 合并音频视频

## 开发环境准备

将 `ffmpeg.exe` 复制到 `cmd` 目录。

## 打包步骤

```shell
cd cmd
go build -o bilidown.exe
```

打包完成后，将 `ffmpeg.exe` 和 `bilidown.exe` 放入压缩包即可。
