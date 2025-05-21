## 开发环境启动指南

- 创建终端（进入 `web` 目录）

    ```bash
    # 启动前端开发服务
    npm run dev
    ```

- 创建终端（进入 `app` 目录）

    ```bash
    # 更新 Electron Forge 配置文件
    tsc forge.config.ts

    # 启动后端 TS 自动编译
    npm run dev:ts
    
    # 启动项目
    npm run start
    ```

## 打包发布

```bash
cd app

# 更新 Electron Forge 配置文件
tsc forge.config.ts

# 生成可执行文件
npm run package

# 生成安装包
npm run make
```

## 源码说明

- `app`：Node.js 环境开发
- `web`：浏览器环境开发