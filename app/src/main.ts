import { BrowserWindow, Menu, app } from 'electron'
import squirrelStartup from 'electron-squirrel-startup'
import { join } from 'path'
import commonIpcMainConfig from './ipcMain/common'
import workIpMainConfig from './ipcMain/work'
import downloadIpMainConfig from './ipcMain/download'
import settingIpMainConfig from './ipcMain/settings'
import { makeMenu } from './menu'
import { openWorkWindow } from './windows'

if (squirrelStartup) app.quit()

commonIpcMainConfig()
workIpMainConfig()
downloadIpMainConfig()
settingIpMainConfig()

app.whenReady().then(() => {
    openWorkWindow()
    makeMenu()
})

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') app.quit()
})