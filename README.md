## bilidown v3

> 哔哩哔哩视频解析下载工具，支持 8K 视频、Hi-Res 音频、杜比视界下载、批量解析，可扫码登录，常驻托盘。

### 技术栈

<table border="1" cellspacing="0" cellpadding="8">
  <thead>
    <tr>
      <th>场景</th>
      <th>名称</th>
      <th>用途</th>
      <th>优点</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="4">前端</td>
      <td><a href="https://vanjs.org/">VanJS</a></td>
      <td>构建响应式 UI</td>
      <td>语法简单，无需复杂环境配置</td>
    </tr>
    <tr>
      <td><a href="https://www.typescriptlang.org/">TypeScript</a></td>
      <td>开发语言</td>
      <td>类型增强，提升可维护性和代码安全</td>
    </tr>
    <tr>
      <td><a href="https://getbootstrap.com/">Bootstrap</a></td>
      <td>样式和组件库</td>
      <td>为原生开发提供优秀的样式和组件</td>
    </tr>
    <tr>
      <td><a href="https://github.com/sindresorhus/p-queue">PQueue</a></td>
      <td>并发操作控制</td>
      <td>控制并发数量，提升并发稳定性</td>
    </tr>
    <tr>
      <td rowspan="4">后端</td>
      <td><a href="https://go.dev/">GoLang</a></td>
      <td>开发语言</td>
      <td>轻松构建跨平台可分发程序</td>
    </tr>
    <tr>
      <td><a href="https://sqlite.org/">SQLite</a></td>
      <td>数据库管理</td>
      <td>无需复杂依赖，配置简单，方便构建</td>
    </tr>
    <tr>
      <td><a href="https://www.ffmpeg.org/">FFmpeg</a></td>
      <td>媒体处理</td>
      <td>强大的媒体处理能力，如音视频合并</td>
    </tr>
    <tr>
      <td><a href="https://github.com/aria2/aria2">Aria2</a></td>
      <td>下载管理</td>
      <td>强大的下载管理能力，支持断点续传</td>
    </tr>
  </tbody>
</table>

### v2 问题收集

- ...

### 技术方案

- 前端通过 `@iuroc/vanjs-utils` 增强样式编程（该库还需优化）
- 前端通过 `@iuroc/vanjs-auto-import-plugin` 实现自动导入（该库还未实现）
- 前端开发支持通过 `@/` 别名访问到 `src/`
- 需要解决 Go 环境调用 ffmpeg 进程弹黑窗口的问题
  - 必要时可以取消托盘，或者尝试写一个启动器，这个启动器启动一个 VBS（仅 Windows），VBS 在后台启动真正的 HTTP 服务，并启用托盘，Linux 平台和 Mac 平台还需要研究一下，v2 中的托盘库在 Linux 测试是没问题的，Mac 上估计也是没问题的
