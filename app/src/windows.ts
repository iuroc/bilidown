import { app, BrowserWindow } from 'electron'
import { join } from 'path'

/** 用于防止重复创建窗口 */
let workWindow: BrowserWindow | null = null
/** 用于防止重复创建窗口 */
let settingsWindow: BrowserWindow | null = null
/** 用于防止重复创建窗口 */
let downloadWindow: BrowserWindow | null = null

const isDev = !app.isPackaged

export function openWorkWindow() {
    if (workWindow && !workWindow.isDestroyed()) {
        workWindow.focus()
        // 窗口可能最小化了，尝试恢复窗口显示
        workWindow.restore()
        return
    }
    workWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../js/preloads/work.js')
        },
        width: 1280,
        height: 800,
        show: false,
    })
    if (isDev) {
        workWindow.loadURL('http://localhost:5173/src/windows/work/index.html')
    } else {
        workWindow.loadFile(join(__dirname, '../dist-renderer/src/windows/work/index.html'))
    }
    workWindow.on('ready-to-show', () => {
        workWindow!.show()
    })
    workWindow.on('closed', () => {
        workWindow = null
    })
}

export function openSettingsWindow() {
    if (settingsWindow && !settingsWindow.isDestroyed()) {
        settingsWindow.focus()
        // 窗口可能最小化了，尝试恢复窗口显示
        settingsWindow.restore()
        return
    }
    settingsWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../js/preloads/settings.js')
        },
        show: false,
        autoHideMenuBar: true
    })
    if (isDev) {
        settingsWindow.loadURL('http://localhost:5173/src/windows/settings/index.html')
    } else {
        settingsWindow.loadFile(join(__dirname, '../dist-renderer/src/windows/settings/index.html'))
    }
    settingsWindow.on('ready-to-show', () => {
        settingsWindow!.show()
    })
    settingsWindow.on('closed', () => {
        settingsWindow = null
    })
}

export function openDownloadWindow() {
    if (downloadWindow && !downloadWindow.isDestroyed()) {
        downloadWindow.focus()
        // 窗口可能最小化了，尝试恢复窗口显示
        downloadWindow.restore()
        return
    }
    downloadWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../js/preloads/download.js')
        },
        show: false,
        autoHideMenuBar: true,
    })
    if (isDev) {
        downloadWindow.loadURL('http://localhost:5173/src/windows/download/index.html')
    } else {
        downloadWindow.loadFile(join(__dirname, '../dist-renderer/src/windows/download/index.html'))
    }
    downloadWindow.on('ready-to-show', () => {
        downloadWindow!.show()
    })
    downloadWindow.on('closed', () => {
        downloadWindow = null
    })
}