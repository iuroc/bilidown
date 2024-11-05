# 开发文档

## 开发环境搭建

### 1. 安装依赖

```bash
# 后端依赖
cd server
go mod download

# 前端依赖
cd ../client
pnpm install
```

### 2. 配置 FFmpeg

将 `ffmpeg.exe` 放入 `server/bin` 目录。构建时，它将被包含在发布包中。

### 3. 启动开发服务器

```bash
# 启动前端开发服务器
cd client
pnpm dev

# 启动后端程序
cd ../server
go run main.go
```

## 构建发布

### 1. 构建发布版本

```bash
# 构建前端
cd client
pnpm build

# 构建并发布后端
cd ../server
goreleaser release --clean
```

### 2. 仅构建不发布

```bash
goreleaser release --snapshot --clean
```

## 待优化部分

-   [ ] 文件名应包含 `bvid`，同 `bvid` 多 `cid` 的应包含序号
