import { State } from 'vanjs-core'
import { goto } from 'vanjs-router'
import { getSeasonInfo, getVideoInfo } from './data'
import { VideoInfoCardData, VideoInfoCardMode } from './type'
import { showErrorPage } from '../mixin'

/** 点击按钮开始解析 */
export const start = async (option: {
    urlInvalid: State<boolean>
    videoInfocardData: State<VideoInfoCardData>
    btnLoading: State<boolean>
    videoInfoCardMode: VideoInfoCardMode
    /** 标识字段类型 */
    idType: IDType
    /** 标识字段值 */
    value: string | number
    /** 调用来源 */
    from: 'click' | 'onfirst'
}) => {
    history.replaceState(null, '', `#/work/${option.idType}/${option.value}`)
    option.btnLoading.val = true
    if (option.idType === 'bv') {
        const bvid = option.value as string
        await getVideoInfo(bvid).then(info => {
            option.videoInfocardData.val = {
                targetURL: `https://www.bilibili.com/video/${bvid}`,
                areas: [],
                styles: [],
                status: '',
                cover: info.pic,
                title: info.title,
                description: info.desc,
                publishData: new Date(info.pubdate * 1000).toLocaleString(),
                duration: info.duration,
                pages: info.pages.map(page => ({ ...page, bvid })),
                dimension: info.dimension,
                owner: info.owner,
                staff: info.staff?.map(i => `${i.name}[${i.title}]`) || []
            }
            option.videoInfoCardMode.val = 'video'
        })
    } else if (option.idType === 'ep' || option.idType === 'ss') {
        const epid = option.idType === 'ep' ? option.value as number : 0
        const ssid = option.idType === 'ss' ? option.value as number : 0
        await getSeasonInfo(epid, ssid).then(info => {
            option.videoInfocardData.val = {
                targetURL: `https://www.bilibili.com/bangumi/play/${option.idType}${option.value}`,
                areas: info.areas.map(i => i.name),
                styles: info.styles,
                duration: 0,
                cover: info.cover,
                description: info.evaluate,
                owner: { name: info.actors, face: '', mid: 0 },
                pages: info.episodes.map((e, index) => ({
                    bvid: e.bvid,
                    cid: e.cid,
                    dimension: e.dimension,
                    duration: e.duration,
                    page: index + 1,
                    part: e.long_title,
                })),
                status: info.new_ep.desc,
                publishData: new Date().toLocaleDateString(),
                staff: info.actors.split('\n'),
                dimension: { height: 0, rotate: 0, width: 0 },
                title: info.title
            }
            option.videoInfoCardMode.val = 'season'
        })
    }
}

export type IDType = 'bv' | 'ep' | 'ss'

/**
 * 校验用户输入的待解析的视频链接
 * @param url 待解析的视频链接
 * @returns 如果校验成功，则返回 BV 号，否则返回 `false`
 */
export const checkURL = (url: string): {
    type: IDType
    value: string | number
} => {
    const matchBvid = url.match(/^(?:https?:\/\/www\.bilibili\.com\/video\/)?(BV1[a-zA-Z0-9]+)/)
    if (matchBvid) return { type: 'bv', value: matchBvid[1] }

    const matchSeason = url.match(/^(?:https?:\/\/www\.bilibili\.com\/bangumi\/play\/)?(ep|ss)(\d+)/)
    if (matchSeason) return { type: matchSeason[1] as 'ep' | 'ss', value: parseInt(matchSeason[2]) }

    throw new Error('您输入的视频链接格式错误')
}

/** 将秒数转换为 `mm:ss` */
export const secondToTime = (second: number) => {
    return `${Math.floor(second / 60)}:${(second % 60).toString().padStart(2, '0')}`
}