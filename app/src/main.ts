import { BrowserWindow, app } from 'electron'

app.whenReady().then(() => {
    const mainWindows = new BrowserWindow({
        webPreferences: {
            // preload: ''
        },
    })

    // mainWindows.loadFile('../web/dist/index.html')
    mainWindows.loadURL('http://localhost:5173/')  // 开发环境
})