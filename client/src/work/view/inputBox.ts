import van from 'vanjs-core'
import { goto } from 'vanjs-router'
import { v4 } from 'uuid'
import { checkURL, handleB23, handleSeasonsArchivesList, start } from '../mixin'
import { WorkRoute } from '..'
import { VanComponent } from '../../mixin'

const { button, div, input, label, span } = van.tags

class InputBoxComp implements VanComponent {
    element: HTMLElement
    btnID = v4()

    constructor(public workRoute: WorkRoute) {
        this.element = div(
            div({ class: () => `hstack gap-3 align-items-stretch ${workRoute.urlInvalidClass.val}` },
                div({ class: () => `form-floating flex-fill` },
                    input({
                        class: () => `form-control border-3 ${workRoute.urlInvalidClass.val}`,
                        placeholder: '请输入待解析的视频链接',
                        value: workRoute.urlValue,
                        oninput: event => workRoute.urlValue.val = event.target.value,
                        onkeyup: event => {
                            if (event.key === 'Enter') document.getElementById(this.btnID)?.click()
                        }
                    }),
                    label({ class: 'w-100' }, '请输入视频链接或 BV/EP/SS 号')
                ),
                ParseButton(this, false, this.btnID),
                ParseButton(this, true)
            ),
            div({ class: 'invalid-feedback' }, () => workRoute.urlInvalid.val ? '您输入的视频链接格式错误' : ''),
        )
    }
}

const ParseButton = (parent: InputBoxComp, large: boolean, id: string = '') => {
    const { workRoute } = parent

    return button({
        class: `btn btn-success text-nowrap ${large ? `btn-lg d-none d-md-block` : 'd-md-none'}`,
        async onclick() {
            try {
                workRoute.btnLoading.val = true
                workRoute.urlValue.val = workRoute.urlValue.val.trim()
                try {
                    const handleB23Result = await handleB23(workRoute.urlValue.val)
                    if (handleB23Result) workRoute.urlValue.val = handleB23Result
                } catch (error) {
                    if (error instanceof Error) return alert(error.message)
                }
                try {
                    const handleSeasonsArchivesListResult = await handleSeasonsArchivesList(workRoute.urlValue.val)
                    if (handleSeasonsArchivesListResult) workRoute.urlValue.val = handleSeasonsArchivesListResult
                } catch (error) {
                    if (error instanceof Error) return alert(error.message)
                }
                const { type, value } = checkURL(workRoute.urlValue.val)
                workRoute.urlInvalid.val = false
                start(workRoute, {
                    idType: type,
                    value,
                    from: 'click'
                }).catch(error => {
                    const errorMessage = `获取视频信息失败：${error.message}`
                    alert(errorMessage)
                    goto('work')
                    workRoute.videoInfoCardMode.val = 'hide'
                })
            } catch (error) {
                workRoute.urlInvalid.val = true
            } finally {
                setTimeout(() => {
                    workRoute.btnLoading.val = false
                }, 200)
            }
        },
        id,
        disabled: workRoute.btnLoading
    }, span({ class: 'spinner-border spinner-border-sm me-2', hidden: () => !workRoute.btnLoading.val }),
        () => workRoute.btnLoading.val ? '解析中' : '解析视频'
    )
}

export default (workRoute: WorkRoute) => new InputBoxComp(workRoute).element