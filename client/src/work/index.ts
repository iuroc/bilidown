import van, { State } from 'vanjs-core'
import { Route, goto, nowHash } from 'vanjs-router'
import { VideoInfoCard, VideoInfoCardData } from './view'
import { getVideoInfo } from './data'
import { v4 } from 'uuid'
import { checkLogin, hasLogin } from '../mixin'

const { button, div, input, label, span } = van.tags

export default () => {
    const urlInvalid = van.state(false)
    const urlInvalidStr = van.derive(() => urlInvalid.val ? 'is-invalid' : '')
    const urlValue = van.state('https://www.bilibili.com/video/BV1H2yfYQEnc/')
    const videoInfocardData = van.state<VideoInfoCardData>({
        title: '', desc: '', pic: '', pubdate: 0, duration: 0, bvid: '',
        pages: [], owner: { face: '', mid: 0, name: '' },
        dimension: { width: 0, height: 0, rotate: 0 }
    })
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
                            label('请输入待解析的视频链接')
                        ),
                        button({
                            class: 'btn btn-primary text-nowrap btn-lg',
                            onclick() {
                                const bvid = checkURL(urlValue.val, urlInvalid)
                                if (!bvid) return
                                start({ urlInvalid, bvid, videoInfocardData, btnLoading })
                            },
                            id: btnID,
                            disabled: btnLoading
                        }, span({ class: 'spinner-border spinner-border-sm me-2', hidden: () => !btnLoading.val }),
                            () => btnLoading.val ? '解析中' : '解析视频'
                        )
                    ),
                    div({ class: 'invalid-feedback' }, () => urlInvalid.val ? '您输入的视频链接格式错误' : ''),
                ),
                VideoInfoCard({ data: videoInfocardData })
            )
        },
        delayed: true,
        async onFirst() {
            if (!await checkLogin()) return
            const bvid = this.args[0]
            if (!bvid) return
            if (bvid && !bvid.match(/^BV1[a-zA-Z0-9]+$/)) return goto('work')
            start({ urlInvalid, bvid, videoInfocardData, btnLoading })
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
    bvid: string
    videoInfocardData: State<VideoInfoCardData>
    btnLoading: State<boolean>
}) => {
    history.replaceState(null, '', `#/work/${option.bvid}`)
    option.btnLoading.val = true
    getVideoInfo(option.bvid).then(info => {
        option.videoInfocardData.val = info
    }).catch(error => {
        alert(`获取视频信息失败：${error.message}`)
        goto('work')
    }).finally(() => {
        option.btnLoading.val = false
    })
}

/**
 * 校验用户输入的待解析的视频链接
 * @param url 待解析的视频链接
 * @param invalid 异常状态值
 * @returns 如果校验成功，则返回 BV 号，否则返回 `false`
 */
const checkURL = (url: string, invalid: State<boolean>) => {
    const match = url.match(/^https?:\/\/www\.bilibili\.com\/video\/(BV1[a-zA-Z0-9]+)/)
    invalid.val = !match
    if (!match) return false
    return match[1]
}