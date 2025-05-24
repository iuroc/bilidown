import { contextBridge, ipcRenderer } from 'electron'

/**
 * 只允许下载管理窗口使用
 */
export const downloadAPI = {
    getAppVersion: () => {
        return ipcRenderer.invoke('common-get-app-version') as Promise<`v${string}`>
    }
}

contextBridge.exposeInMainWorld('downloadAPI', downloadAPI) 