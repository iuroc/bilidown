import { ResJSON } from "../mixin"
import { TaskInDB, TaskStatus, VideoFormat } from "../work/type"

let getActiveTaskController: AbortController | undefined

export const getActiveTask = async (): Promise<ActiveTask[] | null> => {
    getActiveTaskController?.abort()
    getActiveTaskController = new AbortController()
    try {
        const res = await fetch('/api/getActiveTask', {
            signal: getActiveTaskController.signal
        }).then(res => res.json()) as ResJSON<ActiveTask[]>
        if (!res.success) throw new Error(res.message)
        return res.data
    } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') return null
        throw error
    }
}

let getTaskListController: AbortController | undefined

export const getTaskList = async (page: number, pageSize: number): Promise<TaskInDB[] | null> => {
    getTaskListController?.abort()
    getTaskListController = new AbortController()
    try {
        const res = await fetch(`/api/getTaskList?page=${page}&pageSize=${pageSize}`, {
            signal: getTaskListController.signal
        }).then(res => res.json()) as ResJSON<TaskInDB[]>
        if (!res.success) throw new Error(res.message)
        return res.data
    } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') return null
        throw error
    }
}

export const showFile = async (path: string) => {
    const res = await fetch(`/api/showFile?filePath=${encodeURIComponent(path)}`).then(res => res.json()) as ResJSON
    if (!res.success) throw new Error(res.message)
}

/** 用于刷新任务实时进度 */
type ActiveTask = {
    bvid: string
    cid: number
    /** 分辨率代码 */
    format: VideoFormat
    /** 视频标题 */
    title: string
    /** 视频发布者 */
    owner: string
    /** 视频封面 */
    cover: string
    /** 任务进度 */
    status: TaskStatus
    /** 文件保存到的目录 */
    folder: string
    /** 任务 ID */
    id: number
    /** 音频文件下载进度 */
    audioProgress: number
    /** 视频文件下载进度 */
    videoProgress: number
    /** 音视频合并进度 */
    mergeProgress: number
    /** 视频时长，秒 */
    duration: number
}

export const deleteTask = async (id: number) => {
    const res = await fetch(`/api/deleteTask?id=${id}`).then(res => res.json()) as ResJSON
    if (!res.success) throw new Error(res.message)
}