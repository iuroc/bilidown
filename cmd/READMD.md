# bilidown-cmd

Bilibili 视频下载工具。

## 开发前期准备

将 `ffmpeg.exe` 复制到 `cmd` 目录，或在环境变量中配置 `ffmpeg.exe`。

## 打包步骤

```shell
go build -o bilidown.exe
```

打包完成后，将 `ffmpeg.exe` 和 `bilidown.exe` 放入空文件夹中，创建压缩包。
