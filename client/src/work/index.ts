import van, { State } from 'vanjs-core'
import { Route, goto, nowHash } from 'vanjs-router'
import { VideoInfoCard } from './view'
import { getSeasonInfo, getVideoInfo } from './data'
import { v4 } from 'uuid'
import { checkLogin, hasLogin } from '../mixin'
import { VideoInfoCardData } from './type'

const { button, div, input, label, span } = van.tags

export default () => {
    const urlInvalid = van.state(false)
    const urlInvalidStr = van.derive(() => urlInvalid.val ? 'is-invalid' : '')
    // const urlValue = van.state('https://www.bilibili.com/video/BV1H2yfYQEnc/')
    const urlValue = van.state('https://www.bilibili.com/bangumi/play/ep775787')
    const videoInfocardData = van.state<VideoInfoCardData>({
        title: '', description: '', cover: '', publishData: '', duration: 0,
        pages: [], owner: { face: '', mid: 0, name: '' },
        dimension: { width: 0, height: 0, rotate: 0 },
        staff: [], status: '', areas: [], styles: []
    })
    /** 标识视频信息卡片应该显示普通视频还是剧集 */
    const videoInfoCardMode = van.state<'video' | 'season'>('video')
    const btnID = v4()

    const btnLoading = van.state(false)

    return Route({
        rule: 'work',
        Loader() {
            return div({ class: 'vstack gap-3' },
                div(
                    div({ class: () => `hstack gap-3 align-items-stretch ${urlInvalidStr.val}` },
                        div({ class: () => `form-floating flex-fill` },
                            input({
                                class: () => `form-control border-3 ${urlInvalidStr.val}`,
                                placeholder: '请输入待解析的视频链接',
                                type: 'url',
                                value: urlValue,
                                oninput: event => urlValue.val = event.target.value,
                                onkeyup: event => {
                                    if (event.key === 'Enter') document.getElementById(btnID)?.click()
                                }
                            }),
                            label('请输入视频链接或 BV/EP/SS 号')
                        ),
                        button({
                            class: 'btn btn-primary text-nowrap btn-lg',
                            onclick() {
                                try {
                                    const { type, value } = checkURL(urlValue.oldVal)
                                    urlInvalid.val = false
                                    start({
                                        urlInvalid,
                                        videoInfocardData,
                                        btnLoading,
                                        videoInfoCardMode,
                                        idType: type,
                                        value
                                    })
                                } catch (error) {
                                    urlInvalid.val = true
                                }
                            },
                            id: btnID,
                            disabled: btnLoading
                        }, span({ class: 'spinner-border spinner-border-sm me-2', hidden: () => !btnLoading.val }),
                            () => btnLoading.val ? '解析中' : '解析视频'
                        )
                    ),
                    div({ class: 'invalid-feedback' }, () => urlInvalid.val ? '您输入的视频链接格式错误' : ''),
                ),
                VideoInfoCard({ data: videoInfocardData, mode: videoInfoCardMode })
            )
        },
        delayed: true,
        async onFirst() {
            if (!await checkLogin()) return
            const idType = this.args[0] as IDType
            const value = this.args[1]
            if (!value) return goto('work')
            if (idType == 'bv' && !value.match(/^BV1[a-zA-Z0-9]+$/)) return goto('work')
            if ((idType == 'ep' || idType == 'ss') && !value.match(/^\d+$/)) return goto('work')
            start({ urlInvalid, videoInfocardData, btnLoading, videoInfoCardMode, idType, value })
        },
        async onLoad() {
            if (!hasLogin.val) return goto('login')
            this.show()
        }
    })
}

/** 点击按钮开始解析 */
const start = (option: {
    urlInvalid: State<boolean>
    videoInfocardData: State<VideoInfoCardData>
    btnLoading: State<boolean>
    videoInfoCardMode: State<'video' | 'season'>
    /** 标识字段类型 */
    idType: IDType
    /** 标识字段值 */
    value: string | number
}) => {
    history.replaceState(null, '', `#/work/${option.idType}/${option.value}`)
    option.btnLoading.val = true
    if (option.idType === 'bv') {
        const bvid = option.value as string
        getVideoInfo(bvid).then(info => {
            option.videoInfocardData.val = {
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
        }).catch(error => {
            alert(`获取视频信息失败：${error.message}`)
            goto('work')
        }).finally(() => {
            option.btnLoading.val = false
        })
    } else if (option.idType === 'ep' || option.idType === 'ss') {
        const epid = option.idType === 'ep' ? option.value as number : 0
        const ssid = option.idType === 'ss' ? option.value as number : 0
        getSeasonInfo(epid, ssid).then(info => {
            option.videoInfocardData.val = {
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
        }).catch(error => {
            alert(`获取视频信息失败：${error.message}`)
            goto('work')
        }).finally(() => {
            option.btnLoading.val = false
        })
    }
}

export type IDType = 'bv' | 'ep' | 'ss'

/**
 * 校验用户输入的待解析的视频链接
 * @param url 待解析的视频链接
 * @returns 如果校验成功，则返回 BV 号，否则返回 `false`
 */
const checkURL = (url: string): {
    type: IDType
    value: string | number
} => {
    const matchBvid = url.match(/^(?:https?:\/\/www\.bilibili\.com\/video\/)?(BV1[a-zA-Z0-9]+)/)
    if (matchBvid) return { type: 'bv', value: matchBvid[1] }

    const matchSeason = url.match(/^(?:https?:\/\/www\.bilibili\.com\/bangumi\/play\/)?(ep|ss)(\d+)/)
    if (matchSeason) return { type: matchSeason[1] as 'ep' | 'ss', value: parseInt(matchSeason[2]) }

    throw new Error('您输入的视频链接格式错误')
}