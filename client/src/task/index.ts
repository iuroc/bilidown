import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN } from '../mixin'

const { div } = van.tags

export default () => Route({
    rule: 'task',
    Loader() {

        return div('任务列表')
    },
    async onFirst() {
        if (!await checkLogin()) return
    },
    onLoad() {
        if (!GLOBAL_HAS_LOGIN.val) return goto('login')
    },
})