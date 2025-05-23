import { BrowserWindow, app } from 'electron'
import squirrelStartup from 'electron-squirrel-startup'

if (squirrelStartup) app.quit()

app.whenReady().then(() => {
    const mainWindows = new BrowserWindow({
        webPreferences: {
            // preload: ''
        },
    })

    mainWindows.loadURL('http://localhost:5173/')  // 开发环境
})