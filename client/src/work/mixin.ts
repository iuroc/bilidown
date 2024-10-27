import { getSeasonInfo, getVideoInfo } from './data'
import { WorkRoute } from '.'

/** 点击按钮开始解析 */
export const start = async (
    workRoute: WorkRoute,
    option: {
        /** 标识字段类型 */
        idType: IDType
        /** 标识字段值 */
        value: string | number
        /** 调用来源 */
        from: 'click' | 'onfirst'
    }) => {
    history.replaceState(null, '', `#/work/${option.idType}/${option.value}`)
    workRoute.btnLoading.val = true
    workRoute.ownerFaceHide.val = true
    if (option.idType === 'bv') {
        const bvid = option.value as string
        await getVideoInfo(bvid).then(info => {
            workRoute.videoInfocardData.val = {
                targetURL: `https://www.bilibili.com/video/${bvid}`,
                areas: [],
                styles: [],
                status: '',
                cover: info.pic,
                title: info.title,
                description: info.desc,
                publishData: new Date(info.pubdate * 1000).toLocaleString(),
                duration: info.duration,
                pages: info.pages.map((page, index) => ({ ...page, bvid, bandge: (index + 1).toString() })),
                dimension: info.dimension,
                owner: info.owner,
                staff: info.staff?.map(i => `${i.name}[${i.title}]`) || []
            }
            workRoute.videoInfoCardMode.val = 'video'
        })
    } else if (option.idType === 'ep' || option.idType === 'ss') {
        const epid = option.idType === 'ep' ? option.value as number : 0
        const ssid = option.idType === 'ss' ? option.value as number : 0
        await getSeasonInfo(epid, ssid).then(info => {
            workRoute.videoInfocardData.val = {
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
                    bandge: e.title
                })),
                status: info.new_ep.desc,
                publishData: new Date().toLocaleDateString(),
                staff: info.actors.split('\n'),
                dimension: { height: 0, rotate: 0, width: 0 },
                title: info.title
            }
            workRoute.videoInfoCardMode.val = 'season'
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