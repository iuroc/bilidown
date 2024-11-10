# 开发文档

## 打包流程

1. 打包前端

    ```shell
    cd client
    pnpm install
    pnpm build
    ```

2. 打包后端

    ```shell
    cd server
    go mod tidy
    export CGO_ENABLED=1
    go build
    ```

3. 整体打包
    1. 创建一个空文件夹
    2. 将 `server/static` 放入其中
    3. 将 `server/bilidown` 可执行文件放入其中
    4. 将文件夹创建压缩包，分享给任何人

## 交叉编译

1. 构建环境：Ubuntu amd64 24.04
2. 安装 [`osxcross`](https://github.com/tpoechtrager/osxcross) 和 [`goreleaser`](https://goreleaser.com/install/)
3. 将 [`ffmpeg.exe`](https://github.com/iuroc/bilidown/releases/download/v1.0.3/ffmpeg.exe) 放入 `server/bin` 目录
4. 执行命令开始构建

    ```shell
    goreleaser release --snapshot --clean
    ```

## 待优化部分

-   [ ] 保存文件名包含 `bvid`，多 `pages` 的视频同时包含 `(index+1)`

## 相关说明

-   Linux 通过 XVFB 创建虚拟显示器，`dbus-daemon --session --fork` 可以解决其中的警告
