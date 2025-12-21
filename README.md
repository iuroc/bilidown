# bilidown v3

> 哔哩哔哩视频解析下载工具，支持 8K 视频、Hi-Res 音频、杜比视界下载、批量解析，可扫码登录，常驻托盘。

## 技术栈

### 前端

| 名称                                              | 用途          | 优点                             |
| ------------------------------------------------- | ------------- | -------------------------------- |
| [VanJS](https://vanjs.org/)                       | 构建响应式 UI | 语法简单，无需复杂环境配置       |
| [TypeScript](https://www.typescriptlang.org/)     | 开发语言      | 类型增强，提升可维护性和代码安全 |
| [Bootstrap](https://getbootstrap.com/)            | 样式和组件库  | 为原生开发提供优秀的样式和组件   |
| [PQueue](https://github.com/sindresorhus/p-queue) | 并发操作控制  | 控制并发数量，提升并发稳定性     |

### 后端

| 名称                                    | 用途       | 优点                             |
| --------------------------------------- | ---------- | -------------------------------- |
| [GoLang](https://go.dev/)               | 开发语言   | 轻松构建跨平台可分发程序         |
| [SQLite](https://sqlite.org/)           | 数据库管理 | 无需复杂依赖，配置简单，方便构建 |
| [FFmpeg](https://www.ffmpeg.org/)       | 媒体处理   | 强大的媒体处理能力，如音视频合并 |
| [Aria2](https://github.com/aria2/aria2) | 下载管理   | 强大的下载管理能力，支持断点续传 |

## v2 问题收集

- ...

## 技术方案

- 前端通过 `@iuroc/vanjs-utils` 增强样式编程（该库还需优化）
- 前端通过 `@iuroc/vanjs-auto-import-plugin` 实现自动导入（该库还未实现）
- 前端开发支持通过 `@/` 别名访问到 `src/`
- 需要解决 Go 环境调用 ffmpeg 进程弹黑窗口的问题
  - 必要时可以取消托盘，或者尝试写一个启动器，这个启动器启动一个 VBS（仅 Windows），VBS 在后台启动真正的 HTTP 服务，并启用托盘，Linux 平台和 Mac 平台还需要研究一下，v2 中的托盘库在 Linux 测试是没问题的，Mac 上估计也是没问题的。

## 开发说明

- 项目使用 VSCode 开发，通过 [Code Spell Checker](https://marketplace.visualstudio.com/items?itemName=streetsidesoftware.code-spell-checker) 拓展进行拼写检查，开发时可以考虑安装此拓展。

## 开发环境

前端执行 pnpm dev 启动 Vite 开发服务器，其中代理后端地址，开发时访问 Vite 服务地址。

后端执行 air 来启动后端环境，代码修改后会自动重启项目，下面的代码来安装 air。

```bash
go install github.com/air-verse/air@latest
```