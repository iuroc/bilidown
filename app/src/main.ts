import { BrowserWindow, app } from 'electron'
import squirrelStartup from 'electron-squirrel-startup'
import { join } from 'path'

if (squirrelStartup) app.quit()

const isDev = !app.isPackaged

app.whenReady().then(() => {
    const workWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../../js/preloads/work.js')
        },
        show: false,
        autoHideMenuBar: true
    })
    workWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/work/index.html' : '../web/src/windows/work/index.html')
    workWindow.on('ready-to-show', () => {
        workWindow.show()

        setTimeout(() => {
            openSettingsWindow(workWindow)
        }, 3000)
    })
})

export function openSettingsWindow(parent: BrowserWindow) {
    const settingsWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../../js/preloads/settings.js')
        },
        parent,
        alwaysOnTop: true,
        modal: true,
        show: false,
        autoHideMenuBar: true
    })
    settingsWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/settings/index.html' : '../web/src/windows/settings/index.html')
    settingsWindow.on('ready-to-show', () => {
        settingsWindow.show()
    })
}

export function openDownloadWindow() {
    const downloadWindow = new BrowserWindow({
        webPreferences: {
            preload: join(__dirname, '../../js/preloads/download.js')
        },
        autoHideMenuBar: true,
        show: false,
    })
    downloadWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/download/index.html' : '../web/src/windows/download/index.html')
    downloadWindow.on('ready-to-show', () => {
        downloadWindow.show()
    })
}