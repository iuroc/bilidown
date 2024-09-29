# Bilidown

> 哔哩哔哩视频解析下载工具

## 功能规划

1. 本地启动 HTTP 服务器，并将程序隐藏于托盘
2. 启动时自动调用默认浏览器打开服务地址
3. 主页面，进入时会校验二维码有效性，解析按钮按下时也会校验二维码有效性，如果二维码失效，则弹出登录窗口
4. 使用 Nuxt 开发，尝试使用新的 UI 框架
5. 主页面，放置地址输入框，右边一个解析按钮，支持解析下面几种链接：
   1. 普通的 BV1 视频：https://www.bilibili.com/video/BV1NfxMedEUy/
   2. 番剧：https://www.bilibili.com/video/BV1LZ4oe1EcR/ => 重定向到 => https://www.bilibili.com/bangumi/play/ep835909
   3. 带分集的 BV1 视频：https://www.bilibili.com/video/BV1KX4y1V7sA/
6. 主页面输入视频地址，点击解析，如果是单集的视频，直接呈现视频信息和下载界面，如果是多集的，就显示多集列表
7. 组件列表
   1. 分集列表
   2. 单个视频详情卡片
