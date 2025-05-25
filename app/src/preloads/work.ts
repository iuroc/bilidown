import { app, contextBridge, ipcRenderer } from 'electron'

/**
 * 只允许主窗口使用
 */
export const workAPI = {
    /**
     *  打开 [下载管理] 窗口
     */
    openDownloadWindow: () => {
        ipcRenderer.send('work-open-download-window')
    },
    /**
     * 打开 [软件设置] 窗口
     */
    openSettingsWindow: () => {
        ipcRenderer.send('work-open-settings-window')
    },
    /**
     * 获取软件版本号
     * @returns 
     */
    getAppVersion: () => {
        return ipcRenderer.invoke('common-get-app-version') as Promise<`v${string}`>
    }
}

contextBridge.exposeInMainWorld('workAPI', workAPI)