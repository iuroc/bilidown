import { Menu, shell } from 'electron'
import { openDownloadWindow, openSettingsWindow } from './windows'

export function makeMenu() {
    const appMenu = Menu.buildFromTemplate([
        {
            label: '账户',
            submenu: [
                {
                    label: '收藏夹管理'
                },
                {
                    label: '退出登录'
                }
            ]
        },
        {
            label: '下载',
            submenu: [
                {
                    label: '下载管理',
                    click: () => {
                        openDownloadWindow()
                    },
                },
                {
                    label: '打开下载目录',
                    click: () => {
                        openDownloadWindow()
                    },
                },
            ]
        },
        {
            label: '设置',
            click: () => {
                openSettingsWindow()
            },
        },
        {
            label: '帮助',
            click: () => {

            },
            role: 'help'
        },
        {
            label: '更多',
            submenu: [
                {
                    label: '项目主页',
                    click: () => {
                        shell.openExternal('https://github.com/iuroc/bilidown')
                    },
                },
                {
                    label: '支持作者',
                    click: () => {

                    },
                },
                {
                    label: '当前版本',
                    role: 'about',
                },
                { type: 'separator' },
                {
                    label: '开发者工具',
                    role: 'toggleDevTools'
                },
                {
                    label: '关闭软件',
                    role: 'close',
                },
            ]
        },
    ])
    Menu.setApplicationMenu(appMenu)
}