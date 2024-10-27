import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN } from '../mixin'
import { SaveFolderSetting } from './view'

const { button, div } = van.tags

export default () => {

    return Route({
        rule: 'setting',
        Loader() {
            return div({ class: 'vstack gap-4' },
                SaveFolderSetting(),
                div(
                    button({ class: 'btn btn-outline-danger' }, '退出登录')
                )
            )
        },
        async onFirst() {
            if (!await checkLogin()) return
        },
        onLoad() {
            if (!GLOBAL_HAS_LOGIN.val) return goto('login')
        },
    })
}