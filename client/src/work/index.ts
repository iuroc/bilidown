import van, { State } from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import VideoInfoCard from './view/videoInfoCard'
import { checkLogin, GLOBAL_HAS_LOGIN, showErrorPage } from '../mixin'
import { VideoParseResult, VideoInfoCardMode, SectionItem } from './type'
import { IDType, start } from './mixin'
import { ParseModalComp } from './view/parseModal'
import InputBox from './view/inputBox'
import { Modal } from 'bootstrap'

const { div } = van.tags

export class WorkRoute {
    element: HTMLDivElement
    /** 输入框内容是否标识为异常 */
    urlInvalid = van.state(false)
    /** 仅作为类名字符串 */
    urlInvalidClass = van.derive(() => this.urlInvalid.val ? 'is-invalid' : '')
    // urlValue = van.state('https://www.bilibili.com/video/BV1H2yfYQEnc/')
    urlValue = van.state('https://www.bilibili.com/bangumi/play/ep775787')
    videoInfocardData = van.state<VideoParseResult>({
        title: '', description: '', cover: '', publishData: '', duration: 0,
        pages: [], owner: { face: '', mid: 0, name: '' },
        dimension: { width: 0, height: 0, rotate: 0 },
        staff: [], status: '', areas: [], styles: [], targetURL: '',
        section: []
    })
    /** 标识视频信息卡片应该显示普通视频还是剧集，值为 `hide` 时隐藏卡片 */
    videoInfoCardMode: VideoInfoCardMode = van.state('hide')
    ownerFaceHide = van.state(true)

    /** 全部选项卡和列表数据 */
    allSection
    /** 当前选项卡的按钮列表 */
    sectionPages
    /** 当前选中的按钮列表 */
    selectedPages
    /** 视频列表批量解析模态框 */
    parseModal: Modal

    parseModalComp: ParseModalComp

    /** 按钮是否处于 `loading` 状态，如果是则按钮设置为 `disabled` */
    btnLoading = van.state(false)

    sectionTabsActiveIndex = van.state(0)

    constructor() {
        const _that = this
        this.allSection = van.derive(() => [{ title: '正片', pages: this.videoInfocardData.val.pages }]
            .concat(this.videoInfocardData.val.section))
        this.sectionPages = van.derive(() => this.allSection.val[this.sectionTabsActiveIndex.val].pages || [])
        this.selectedPages = van.derive(() => this.sectionPages.val.filter(page => page.selected.val))
        this.parseModalComp = new ParseModalComp({ workRoute: this })
        van.add(document.body, this.parseModalComp.element)

        this.parseModal = new Modal(this.parseModalComp.element)

        this.element = Route({
            rule: 'work',
            Loader() {
                return div({ class: 'vstack gap-3' },
                    InputBox(_that),
                    div({ hidden: () => _that.videoInfoCardMode.val == 'hide' || _that.btnLoading.val },
                        VideoInfoCard(_that),
                    ),
                )
            },
            async onFirst() {
                if (!await checkLogin()) return
                const idType = this.args[0] as IDType
                const value = this.args[1]
                if (!value) return goto('work')
                if (idType == 'bv' && !value.match(/^BV1[a-zA-Z0-9]+$/)) return goto('work')
                if ((idType == 'ep' || idType == 'ss') && !value.match(/^\d+$/)) return goto('work')
                if (idType == 'bv') _that.urlValue.val = value
                else if (idType == 'ep' || idType == 'ss') _that.urlValue.val = `${idType}${value}`
                start(_that, {
                    idType,
                    value,
                    from: 'onfirst'
                }).catch(error => {
                    const errorMessage = `获取视频信息失败：${error.message}`
                    showErrorPage(errorMessage)
                    _that.videoInfoCardMode.val = 'hide'
                }).finally(() => {
                    _that.btnLoading.val = false
                })
            },
            async onLoad() {
                if (!GLOBAL_HAS_LOGIN.val) return goto('login')
            }
        })
    }
}

export default () => new WorkRoute().element