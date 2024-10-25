import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { checkLogin, hasLogin } from '../mixin'

const { div } = van.tags

export default () => Route({
    rule: 'task',
    Loader() {

        return div('任务列表')
    },
    delayed: true,
    async onFirst() {
        if (!await checkLogin()) return
    },
    onLoad() {
        if (!hasLogin.val) return goto('login')
        this.show()
    },
})