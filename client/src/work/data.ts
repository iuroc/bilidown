import { ResJSON, timeoutController } from '../mixin'
import { SeasonInfo, VideoInfo } from './type'

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


/*
影视飓风[UP主], XXXX[YYYYY],XXXX[YYYYY]    王大头 赵铁柱 赵煊 苏子缙 陈胖胖 陈皮诺Pino 张耀锋 刘莹莹
1                                         已完结, 全12集
1920x1080                                 中国大陆
时长                                       短剧/搞笑/剧情
*/