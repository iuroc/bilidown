import van from 'vanjs-core'
import { goto } from 'vanjs-router'
import { v4 } from 'uuid'
import { checkURL, start } from '../mixin'
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
                        type: 'url',
                        value: workRoute.urlValue,
                        oninput: event => workRoute.urlValue.val = event.target.value,
                        onkeyup: event => {
                            if (event.key === 'Enter') document.getElementById(this.btnID)?.click()
                        }
                    }),
                    label('请输入视频链接或 BV/EP/SS 号')
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
        onclick() {
            try {
                const { type, value } = checkURL(workRoute.urlValue.oldVal)
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
                }).finally(() => {
                    workRoute.btnLoading.val = false
                })
            } catch (error) {
                workRoute.urlInvalid.val = true
            }
        },
        id,
        disabled: workRoute.btnLoading
    }, span({ class: 'spinner-border spinner-border-sm me-2', hidden: () => !workRoute.btnLoading.val }),
        () => workRoute.btnLoading.val ? '解析中' : '解析视频'
    )
}

export default (workRoute: WorkRoute) => new InputBoxComp(workRoute).element