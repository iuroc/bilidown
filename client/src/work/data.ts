import { ResJSON, timeoutController } from '../mixin'
import { PlayInfo, SeasonInfo, TaskInitData, VideoInfo } from './type'

/**
 * 获取视频信息
 * 
 * @throws {Error}
 */
export const getVideoInfo = async (bvid: string): Promise<VideoInfo> => {
    const { signal, timer } = timeoutController()

    const res = await fetch(`/api/getVideoInfo?bvid=${bvid}`, {
        signal
    }).then(res => res.json()) as ResJSON<VideoInfo>
    if (!res.success) throw new Error(res.message)
    clearTimeout(timer)
    return res.data
}

/**
 * 获取剧集信息
 * @param epid EP 号
 * @param ssid SS 号
 * @throws {Error}
 */
export const getSeasonInfo = async (epid: number, ssid: number): Promise<SeasonInfo> => {
    const { signal, timer } = timeoutController()

    const res = await fetch(`/api/getSeasonInfo?epid=${epid}&ssid=${ssid}`, {
        signal
    }).then(res => res.json()) as ResJSON<SeasonInfo>

    if (!res.success) throw new Error(res.message)
    clearTimeout(timer)
    return res.data
}


export const getPlayInfo = async (bvid: string, cid: number, controller: AbortController): Promise<PlayInfo> => {
    const res = await fetch(`/api/getPlayInfo?bvid=${bvid}&cid=${cid}`, {
        signal: controller.signal
    }).then(res => res.json()) as ResJSON<PlayInfo>
    if (!res.success) throw new Error(res.message)
    return res.data
}

export const createTask = async (tasks: TaskInitData[]): Promise<ResJSON> => {
    const { signal, timer } = timeoutController()

    return await fetch('/api/createTask', {
        method: 'POST',
        signal,
        body: JSON.stringify(tasks),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => res.json()) as ResJSON
}