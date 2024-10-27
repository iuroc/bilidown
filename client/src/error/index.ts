import van from 'vanjs-core'
import { Route } from 'vanjs-router'
import { GLOBAL_ERROR_MESSAGE, GLOBAL_HIDE_PAGE } from '../mixin'

const { button, div } = van.tags

export default () => Route({
    rule: 'error',
    Loader() {
        return div({ class: 'py-5 px-4 container' },
            div({ class: 'py-5 px-3 border rounded-4 vstack align-items-center gap-4' },
                div({ class: 'fs-2 fw-bold' }, '错误提示'),
                div({ class: 'text-danger fs-4' }, () => GLOBAL_ERROR_MESSAGE.val || '请刷新页面重试'),
                button({
                    class: 'btn btn-warning', onclick() {
                        location.href = location.pathname
                    }
                }, '刷新页面')
            )
        )
    },
})