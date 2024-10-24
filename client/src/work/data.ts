import { ResJSON, timeoutController } from '../mixin'
import { VideoInfo } from './view'

export const getVideoInfo = async (bvid: string): Promise<VideoInfo> => {
    const { signal, timer } = timeoutController()

    const res = await fetch(`/api/getVideoInfo?bvid=${bvid}`, {
        signal
    }).then(res => res.json()) as ResJSON<VideoInfo>
    if (!res.success) throw new Error(res.message)
    clearTimeout(timer)
    return res.data
}
