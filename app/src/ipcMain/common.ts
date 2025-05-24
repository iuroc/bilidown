import { app, ipcMain } from 'electron'
import { clean } from 'semver'

export default () => {
    ipcMain.handle('common-get-app-version', () => {
        return `v${clean(app.getVersion())}`
    })
}