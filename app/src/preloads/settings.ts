import { contextBridge, ipcRenderer } from 'electron'

/**
 * 只允许软件设置窗口使用
 */
export const settingsAPI = {
    getAppVersion: () => {
        return ipcRenderer.invoke('common-get-app-version') as Promise<`v${string}`>
    }
}

contextBridge.exposeInMainWorld('settingsAPI', settingsAPI)