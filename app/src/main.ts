import { BrowserWindow, app } from 'electron'
import squirrelStartup from 'electron-squirrel-startup'

if (squirrelStartup) app.quit()

const isDev = !app.isPackaged

app.whenReady().then(() => {
    const workWindow = new BrowserWindow({
        webPreferences: {
            // preload: ''
        },
        autoHideMenuBar: true
    })

    const downloadWindow = new BrowserWindow({
        webPreferences: {
            // preload: ''
        },
        autoHideMenuBar: true
    })

    const settingsWindow = new BrowserWindow({
        webPreferences: {
            // preload: ''
        },
        autoHideMenuBar: true
    })

    workWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/work/index.html' : '../web/src/windows/work/index.html')
    downloadWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/download/index.html' : '../web/src/windows/download/index.html')
    settingsWindow.loadURL(isDev ? 'http://localhost:5173/src/windows/settings/index.html' : '../web/src/windows/settings/index.html')
})