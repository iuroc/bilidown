import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { VideoInfoCard } from './view'
import { v4 } from 'uuid'
import { checkLogin, GLOBAL_HAS_LOGIN, showErrorPage } from '../mixin'
import { VideoInfoCardData, VideoInfoCardMode } from './type'
import { IDType, checkURL, start } from './mixin'

const { button, div, input, label, span } = van.tags

export default () => {
    const urlInvalid = van.state(false)
    const urlInvalidStr = van.derive(() => urlInvalid.val ? 'is-invalid' : '')
    const urlValue = van.state('https://www.bilibili.com/video/BV1H2yfYQEnc/')
    // const urlValue = van.state('https://www.bilibili.com/bangumi/play/ep775787')
    const videoInfocardData = van.state<VideoInfoCardData>({
        title: '', description: '', cover: '', publishData: '', duration: 0,
        pages: [], owner: { face: '', mid: 0, name: '' },
        dimension: { width: 0, height: 0, rotate: 0 },
        staff: [], status: '', areas: [], styles: [], targetURL: ''
    })
    /** 标识视频信息卡片应该显示普通视频还是剧集，值为 `hide` 时隐藏卡片 */
    const videoInfoCardMode: VideoInfoCardMode = van.state('hide')
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
                                        value,
                                        from: 'click'
                                    }).catch(error => {
                                        const errorMessage = `获取视频信息失败：${error.message}`
                                        alert(errorMessage)
                                        goto('work')
                                        videoInfoCardMode.val = 'hide'
                                    }).finally(() => {
                                        btnLoading.val = false
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
        async onFirst() {
            if (!await checkLogin()) return
            const idType = this.args[0] as IDType
            const value = this.args[1]
            if (!value) return goto('work')
            if (idType == 'bv' && !value.match(/^BV1[a-zA-Z0-9]+$/)) return goto('work')
            if ((idType == 'ep' || idType == 'ss') && !value.match(/^\d+$/)) return goto('work')
            if (idType == 'bv') urlValue.val = value
            else if (idType == 'ep' || idType == 'ss') urlValue.val = `${idType}${value}`
            start({ urlInvalid, videoInfocardData, btnLoading, videoInfoCardMode, idType, value, from: 'onfirst' })
                .catch(error => {
                    const errorMessage = `获取视频信息失败：${error.message}`
                    showErrorPage(errorMessage)
                    videoInfoCardMode.val = 'hide'
                }).finally(() => {
                    btnLoading.val = false
                })
        },
        async onLoad() {
            if (!GLOBAL_HAS_LOGIN.val) return goto('login')
        }
    })
}
