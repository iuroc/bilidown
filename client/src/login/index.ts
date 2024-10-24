import van from 'vanjs-core'
import { Route } from 'vanjs-router'

const { div } = van.tags

export default () => {

    return Route({
        rule: 'login',
        Loader() {

            return div('登录页面')
        },
    })
}