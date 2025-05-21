import van from 'vanjs-core'
import { Route } from 'vanjs-router'

const { div, button, a: anchor } = van.tags

export default () => {
    return Route({
        rule: 'home',
        Loader() {

            return div(
                button({
                    onclick: () => {
                        alert('Good!')
                    }
                }, 'Click Me')
            )
        },
    })
}