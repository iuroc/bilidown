# Bilidown

> 哔哩哔哩视频解析下载工具

## 开发环境

1. 将 `gcc` 命令配置到 `PATH` 环境变量
2. 设置环境变量 `CGO_ENABLED=1`
3. 执行下面的代码：

    ```shell
    go install github.com/air-verse/air@latest
    npm install pnpm -g

    cd server
    go mod tidy

    cd ../client
    pnpm install
    pnpm dev
    ```

## 编译运行

```shell
pnpm build
pnpm start
```

## 功能规划

1. 本地启动 HTTP 服务器，并将程序隐藏于托盘
2. 启动时自动调用默认浏览器打开服务地址
3. 主页面，进入时会校验二维码有效性，解析按钮按下时也会校验二维码有效性，如果二维码失效，则弹出登录窗口
4. 使用 Nuxt 开发，尝试使用新的 UI 框架
5. 主页面，放置地址输入框，右边一个解析按钮，支持解析下面几种链接：
    1. 普通的 BV1 视频：https://www.bilibili.com/video/BV1NfxMedEUy/
    2. HDR 视频：https://www.bilibili.com/video/BV1rp4y1e745/
    3. 番剧：https://www.bilibili.com/video/BV1LZ4oe1EcR/ => 重定向到 => https://www.bilibili.com/bangumi/play/ep835909
    4. 带分集的 BV1 视频：https://www.bilibili.com/video/BV1KX4y1V7sA/
6. 主页面输入视频地址，点击解析，如果是单集的视频，直接呈现视频信息和下载界面，如果是多集的，就显示多集列表
7. 组件列表
    1. 分集列表
    2. 单个视频详情卡片

## 关键接口

1. 根据 bvid 获取视频信息（包括基本信息和分集列表，含 bvid 和 cid）：https://api.bilibili.com/x/web-interface/view?bvid=（携带 Cookie）
2. 根据 bvid 和 cid 获取视频播放地址：https://api.bilibili.com/x/player/wbi/playurl?bvid=&cid=&fnval=4048&fnver=0&fourk=1
3. 根据 epid 获取番剧信息（包括基本信息和分集列表，含 bvid 和 cid）：https://api.bilibili.com/pgc/view/web/season?ep_id=

## 功能规划（越往下越新）

### 视频解析（/work）

输入框和解析按钮，点击后下面呈现分集列表，如果是只有 1 项的，那就只显示 1 项，点击分集列表项，可以弹出模态框查看视频的详细信息，不管是当个还是批量，都默认全部选中分集结果，然后底部是【打包解析下载】按钮（单个结果时则显示【解析下载】），然后将下载任务以 JSON POST 的方式发给后端，后端向服务器创建一个下载任务，前端会询问（弹出窗口，选择下载到哪个文件夹），然后前端弹出提示，下载任务创建成功，创建任务时，后端向数据库插入任务记录，比如是批量解析的，其实不需要合并为一个组，只需要在批量解析下载时，选择一个统一的目录即可，而存储下载记录时，依然是分开来单个存储的（记录包括：bvid, cid, time 时间, status 进度, path 完整路径，视频标题、封面、下载时间、文件大小、分辨率、UP 主名称），后端 Go 在每次启动 HTTP 服务器时，都将未完成的数据库记录的进度改为失败，HTTP 服务器存活期间，Go 的临时的 Sync Map 数据是持续有效的，适合存储每个任务的实时下载进度，并供前端轮询。前端轮询时，Go 将 Sync Map 中未完成的（已经完成的记录应该等数据库进度更新为完成后就从 Map 删除）记录包括进度一起返回给前端。

### 下载任务（/task）

展示历史下载记录（包括视频标题、封面、下载时间、文件大小、分辨率、UP 主名称）

参考文档：https://socialsisteryi.github.io/bilibili-API-collect/docs/video/videostream_url.html

校验是否登录（携带 Cookie）：https://api.bilibili.com/x/space/myinfo

## 功能规划

SESSDATA 存储在 SQLite，前端每次初始进入页面时，都从 SQLite 请求一次 SESSDATA，后端校验 SESSDATA 有效性，如果无效，则通知前端获取登陆二维码进行登录，扫码登录完成后将 SESSDATA 汇报给后端，前端环境不存储 SESSDATA。

## 数据库结构 `data.db`

### field 字段表

-   `key` string
-   `value` string

该表插入 SESSDATA.

### task 任务表

-   `title` 视频标题，属于当前 cid
-   `bvid`
-   `cid`
-   `cover` 视频封面
-   `desc` 视频简介
-   `owner` UP 主名称
-   `owner_cover` UP 主头像
-   `create_time` 任务创建时间
-   `duration` 视频时常
-   `path` 保存的绝对路径
-   `status` 任务进度，done, waiting, running, error

## 关于下载

输入视频地址，点击解析按钮，显示详情卡片，如果有多集，底下出现分集标题列表。

如果没有多集的，则同时获取视频的支持格式列表，下拉菜单选择后，底下有一个输入框，里面是保存路径，再下面是下载按钮“点击下载”

保存路径中会显示默认的保存路径，如果路径为空，则下载按钮为灰色。

可以点击边上按钮选择保存的文件夹。

如果是多集的，显示分集标题列表，每项左边都有复选框，底下按钮是“解析选中项目”，点击后，
按钮变为 loading 状态，并显示 m/n 进度，全部完成后，弹出模态框，是一个列表，里面是刚才选中的全部视频标题，

以及标题右边是一个下拉框，里面包含了该分集视频支持的格式列表，列表的顶部支持为全部视频选择格式，这个下拉菜单的选项，
则是从全部视频的合并结果来的，有效分集不一定支持，默认情况下，每一项都默认选中本身支持的最高质量，而如果是被顶部下拉菜单
统一设置，则当分集本身不支持该格式时，选一个最近的格式进行降级。

模态框底部是“开始下载”，会将全部视频的 bvid、cid、期望的格式 ID 发送给后端进行任务创建。

## 等待更新

-   移动端的分享链接解析

## 关于文件命名

1. `[序号] 主标题 [分集标题，如果同主标题则省略] [发布者] [分辨率] [时长] 随机码`
2. `长标题 [短标题] [主标题] [分辨率] [时长] 随机码`

## 代码优化说明

1. 把不必要的大写开头导出取消
2. 保证 log.Fatal 只能在入口层调用，其余层都使用 err 返回值
3. 适当调整函数封装，完善和修正注释内容
4. 集中构建错误响应信息，如：

    ```go
    if r.Method != http.MethodPost {
    	util.Res{Success: false, Message: "不支持的请求方法"}.Write(w)
    	return
    }
    ```

    改成：

    ```go
    if r.Method != http.MethodPost {
    	errutil.write(w, errutil.MethodNotAllow)
    	return
    }
    ```

    ```go
    package errutil

    const (
        MethodNotAllow = "不支持的请求方法"
        ParamsError = "参数错误"
    )
    ```

5. 还是不太能区分为 type 增加方法时，对 type 的引用应该用指针还是不用的区别。
6. 前端增加一些 Loading 动画
7. 注意非 Windows 的构建不要包含 exe 文件，这个看能不能在 goreleaser 里去配置，然后 Windows 版本也是可以生成 2 份，分别是带 ffmpeg 和不带的
8. 支持手机分享链接解析

## 容器操作

```shell
docker run -it -v data:/usr/src/data -w /usr/src/data -d -p 8100:8098 golang
```

```shell
apt update
apt install -y libayatana-appindicator3-1  # go build 编译需要
apt install -y ffmpeg  # 运行时需要
apt install -y vim  # 开发工具，可删除
apt install -y xvfb  # 虚拟显示器，运行时需要
apt install -y dbus-x11  # 运行时需要
wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -  # 安装 pnpm
source $HOME/.bashrc  # 使 pnpm 生效
apt install -y nodejs  # 前端编译需要
apt install -y libayatana-appindicator3-dev  # go build 编译需要
```

```shell
# 运行程序
Xvfb :99 -screen 0 1920x1080x24 &
export DISPLAY=:99
eval $(dbus-launch --sh-syntax)
export DBUS_SESSION_BUS_ADDRESS
./bilidown
```

```shell
cd /usr/src
git clone https://github.com/tpoechtrager/osxcross
cd osxcross
wget -P tarballs https://github.com/alexey-lysiuk/macos-sdk/releases/download/14.5/MacOSX14.5.tar.xz
apt install clang cmake -y
apt install -y libssl-dev
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
apt update
apt install goreleaser
```

```shell
docker pull iuroc/cgo-cross-build:latest
docker run --rm -v .:/usr/src/data cgo-cross-build goreleaser release --snapshot --clean
```
