@README.md

# AGENTS.md — 开发者指南

本文件面向开发者与 AI 编程代理（Claude Code 通过根目录的 CLAUDE.md `@AGENTS.md` 自动加载本文件）。用户向的介绍见上方导入的 README.md，此处只讲开发相关的内容。

## 项目概览

Bilidown 是「本地 Web 应用 + 系统托盘」架构的单机工具：

- `server/main.go` 启动后常驻系统托盘（[getlantern/systray](https://github.com/getlantern/systray)），并在本机 **8098 端口**（硬编码于 `main.go` 常量）启动 HTTP 服务，随后自动调用默认浏览器打开 `http://127.0.0.1:8098`。
- 前端是单页应用（Vite + VanJS + Bootstrap + TypeScript），构建产物输出到 `server/static`，由 Go 进程以 `http.Dir("static")` 直接托管，**运行时依赖该目录存在**。
- 数据存储使用 SQLite（纯 Go 驱动 `modernc.org/sqlite`，无 CGO 依赖），数据库文件与下载的视频均位于程序工作目录（默认下载目录 `./download`，可在前端设置页修改）。
- 下载依赖 FFmpeg 合并音视频：程序在 PATH 中查找 `ffmpeg`，找不到再查找程序目录下 `bin/ffmpeg`（见 `server/util/util.go` 的 `GetFFmpegPath`）；FFmpeg 缺失时程序打印提示并拒绝启动。

## 目录结构

```text
client/                  前端源码（Vite + VanJS + Bootstrap + TypeScript）
  src/                   前端源代码（work 解析下载、login 扫码登录、setting 设置等模块）
  public/locales/        前端文案
  vite.config.ts         构建配置，产物输出到 ../server/static
server/                  Go 后端
  main.go                入口：托盘、HTTP 服务、数据库初始化与迁移
  bilibili/              B 站 API 封装
  router/                /api 路由（video 解析、task 任务、login 登录）
  task/                  下载任务调度
  util/                  通用工具（DB、FFmpeg 查找、统一 JSON 响应等）
  static/                前端构建产物（gitignore，由 pnpm build 生成）
  bin/                   本地放置的 FFmpeg（gitignore）
.github/workflows/       release.yml：GitHub Actions 自动构建发布
```

## 本地开发

前置要求：Go 1.23+、Node.js + pnpm 9+、FFmpeg（在 PATH 或 `server/bin/` 中）。Linux 还需要 GTK 开发文件：

```shell
sudo apt install pkg-config gcc libgtk-3-dev libayatana-appindicator3-dev
```

开发流程：

```shell
# 1. 前端（server/static 已存在时可跳过构建，直接改后端调试）
cd client
pnpm install
pnpm build      # 产物输出到 ../server/static

# 2. 后端
cd server
go build && ./bilidown    # 或 go run main.go
```

前端开发时可在 `client/` 目录用 `pnpm dev` 启动 Vite 热更新；dev 服务器已把 `/api` 代理到 `http://127.0.0.1:8098`（见 `vite.config.ts`），因此调试接口时需要后端同时在 8098 端口运行。

前后端联调时，在仓库根目录运行 `pnpm dev`：它并行启动前端 Vite 热更新与后端 air 热重载（`pnpm -r --parallel`）。air 需单独安装：`go install github.com/air-verse/air@latest`。`pnpm build` 后可运行 `pnpm start` 直接启动打包好的 `bilidown` 可执行文件。

日常只需改后端时，在 `server/` 目录直接运行 `main.go` 即可，无需任何交叉编译或容器。

## 构建

```shell
cd client && pnpm build    # 先构建前端到 server/static
cd ../server && go build   # 再构建后端
```

- Windows：systray 走纯 Go 的 Win32 API，可用 `CGO_ENABLED=0 go build`；正式发行加 `-ldflags "-H windowsgui"` 隐藏控制台窗口。
- Linux / macOS：需要 CGO（GTK / AppKit），保持默认 `CGO_ENABLED=1` 直接 `go build`。

**本仓库不提供交叉编译。** 如需其他平台的产物：推送 `v*` 标签由 GitHub Actions 构建（见下），或在对应操作系统上自行执行上述构建命令。

## 发布流程（GitHub Actions）

推送 `v*` 格式的标签（如 `v2.1.2`）即触发 [.github/workflows/release.yml](.github/workflows/release.yml)：

1. `frontend` 任务构建一次前端，产物 `server/static` 通过 artifact 共享。
2. `build` 任务在各目标平台的**原生 runner** 上构建后端：Linux x86_64、Windows x86_64/i386/arm64、macOS arm64（Intel Mac 配置暂时注释）。Windows 包在 CI 中从本仓库 v1.0.1 release 下载 `ffmpeg.exe` 一并打入。
3. `release` 任务汇总产物并自动创建 GitHub Release。

产物命名：`bilidown_Linux_x86_64.tar.gz`、`bilidown_Windows_x86_64.zip` / `bilidown_Windows_i386.zip` / `bilidown_Windows_arm64.zip`、`bilidown_Darwin_arm64.tar.gz`。也可在 Actions 页面手动触发（workflow_dispatch）做测试构建，不会发布。

## 代码约定

- 注释与 UI 文案使用中文。
- API 响应统一使用 `util.Res{Success, Message, Data}` JSON 结构（见 `server/util/response.go`）。
- 访问 SQLite 前必须持有 `util.SqliteLock` 锁。
- `Must*` 前缀的函数表示失败直接 `log.Fatalln` 退出（如 `mustInitTables`）。
- 数据库结构变更通过 `main.go` 的 `addMissingColumns` 做增量迁移（SQLite 加列兼容旧库）。
