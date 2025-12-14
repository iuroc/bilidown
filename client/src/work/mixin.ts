import { getFavList, getRedirectedLocation, getSeasonInfo, getVideoInfo } from './data'
import { WorkRoute } from '.'
import van from 'vanjs-core'
import { Episode, PageInParseResult, VideoParseResult } from './type'
import { ResJSON } from '../mixin'

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
    // 如果初始载入路由时，路由参数不带有视频信息，则会随机选择一个热门视频
    // 这种情况展示的解析结果，不应该同步修改路由参数，用户直接刷新网页，可以继续随机推荐不同的视频
    // 而如果是单击按钮或者路由参数带有视频信息，那么实现的解析结果需要同步修改路由参数
    if (!workRoute.isInitPopular.val)
        history.replaceState(null, '', `#/work/${option.idType}/${option.value}`)
    workRoute.sectionTabsActiveIndex.val = 0
    if (option.idType === 'bv') {
        const bvid = option.value as string
        await getVideoInfo(bvid).then(info => {
            workRoute.videoInfocardData.val = {
                section: (info.ugc_season.sections || []).map(i => {
                    let index = 0
                    return {
                        title: (info.ugc_season.sections || []).length == 1 ? info.ugc_season.title : i.title,
                        pages: i.episodes.flatMap(j => j.pages.map(k => {
                            index++
                            return {
                                bvid: j.bvid,
                                cid: k.cid,
                                dimension: k.dimension,
                                duration: k.duration,
                                page: index,
                                part: k.part || j.title,  // 优先使用page的part字段，如果为空则使用episode的title
                                bandge: index.toString(),
                                selected: van.state(bvid == j.bvid)
                            }
                        }))
                    }
                }),
                targetURL: `https://www.bilibili.com/video/${bvid}`,
                areas: [],
                styles: [],
                status: '',
                cover: info.pic,
                title: info.title,
                description: info.desc,
                publishData: new Date(info.pubdate * 1000).toLocaleString(),
                duration: info.duration,
                pages: info.pages.map((page, index) => ({
                    ...page,
                    bvid,
                    bandge: (index + 1).toString(),
                    selected: van.state(info.pages.length == 1)
                })),
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
                section: (info.section || []).map(i => ({
                    pages: i.episodes.map(episodeToPage),
                    title: i.title
                })),
                targetURL: `https://www.bilibili.com/bangumi/play/${option.idType}${option.value}`,
                areas: info.areas.map(i => i.name),
                styles: info.styles,
                duration: 0,
                cover: info.cover,
                description: info.evaluate,
                owner: { name: info.actors, face: '', mid: 0 },
                pages: info.episodes.map(episodeToPage),
                status: info.new_ep.desc,
                publishData: new Date().toLocaleDateString(),
                staff: info.actors.split('\n'),
                dimension: { height: 0, rotate: 0, width: 0 },
                title: info.title
            }
            workRoute.videoInfoCardMode.val = 'season'
        })
    } else if (option.idType == 'fav') {
        const mediaId = option.value as number
        await getFavList(mediaId).then(favList => {
            workRoute.videoInfocardData.val = {
                areas: [],
                cover: favList[0].cover,
                description: favList[0].intro,
                dimension: { height: 0, width: 0, rotate: 0 },
                duration: favList[0].duration,
                owner: favList[0].upper,
                pages: favList.map((info, index) => ({
                    bandge: (index + 1).toString(),
                    selected: van.state(index == 0),
                    bvid: info.bvid,
                    cid: info.ugc.first_cid,
                    dimension: { height: 0, width: 0, rotate: 0 },
                    duration: info.duration,
                    page: index,
                    part: info.title,
                })),
                publishData: new Date(favList[0].pubtime * 1000).toLocaleString(),
                section: [],
                staff: [],
                status: '',
                styles: [],
                targetURL: `https://www.bilibili.com/video/${favList[0].bvid}`,
                title: favList[0].title
            }
            workRoute.videoInfoCardMode.val = 'video'
        })
    }
}

const episodeToPage = (episode: Episode, index: number): PageInParseResult => {
    return {
        bvid: episode.bvid,
        cid: episode.cid,
        dimension: episode.dimension,
        duration: episode.duration,
        page: index + 1,
        part: episode.long_title || (episode.title.match(/^\d+$/) ? `第 ${episode.title} 集` : episode.title),
        bandge: episode.title,
        selected: van.state(false)
    }
}

export type IDType = 'bv' | 'ep' | 'ss' | 'fav'

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

    try {
        const _url = new URL(url)
        const mediaId = parseInt(_url.searchParams.get('fid') || '')
        if (_url.hostname == 'space.bilibili.com' && _url.pathname.match(/^\/\d+\/favlist$/) && !isNaN(mediaId)) {
            return { type: 'fav', value: mediaId }
        }
        const mlMatch = url.match(/^https:\/\/www.bilibili.com\/medialist\/detail\/ml(\d+)/)
        if (mlMatch) return { type: 'fav', value: mlMatch[1] }
    } catch { }
    throw new Error('您输入的视频链接格式错误')
}

/** 将秒数转换为 `mm:ss` */
export const secondToTime = (second: number) => {
    return `${Math.floor(second / 60)}:${(second % 60).toString().padStart(2, '0')}`
}

/** 如果是 B23 地址，则返回重定向后的地址，否则返回 `false` */
export const handleB23 = async (url: string): Promise<string | false> => {
    if (!url.match(/^https:\/\/b23.tv\//)) return false
    const epMatch = url.match(/^https:\/\/b23.tv\/(ep|ss)(\d+)/)
    if (epMatch) return `https://www.bilibili.com/bangumi/play/${epMatch[1]}${epMatch[2]}`
    const location = await getRedirectedLocation(url)
    return location
}

export const handleSeasonsArchivesList = async (url: string): Promise<string | false> => {
    try { new URL(url) } catch { return false }
    const _url = new URL(url)
    const mid = _url.pathname.match(/^\/(\d+)\/channel\/collectiondetail$/)?.[1]
    const seasonId = parseInt(_url.searchParams.get('sid') || '')
    if (_url.hostname == 'space.bilibili.com' && mid && !isNaN(seasonId)) {
        return fetch(`/api/getSeasonsArchivesListFirstBvid?mid=${mid}&seasonId=${seasonId}`)
            .then(res => res.json())
            .then((body: ResJSON<string>) => {
                if (!body.success) throw new Error(body.message)
                return body.data
            })
    }
    return false
}