import van from 'vanjs-core'
import { Route } from 'vanjs-router'

const { div } = van.tags

export default () => Route({
    rule: 'task',
    Loader() {

        return div('任务列表')
    },
})