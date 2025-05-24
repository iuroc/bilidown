import { ipcMain } from 'electron'
import commonIpcMain from './common'
import { openDownloadWindow, openSettingsWindow } from '../windows'

export default () => {
    ipcMain.on('work-open-download-window', () => {
        openDownloadWindow()
    })

    ipcMain.on('work-open-settings-window', () => {
        openSettingsWindow()
    })
}