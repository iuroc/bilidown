import van, { State } from 'vanjs-core'
import { Route, goto, now } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN, GLOBAL_HIDE_PAGE, ResJSON, VanComponent } from '../mixin'
import { getActiveTask, getTaskList } from './data'
import { TaskInDB, TaskStatus } from '../work/type'

const { div } = van.tags

export class TaskRoute implements VanComponent {
    element: HTMLElement

    taskList: State<(TaskInDB & {
        /** 音频下载进度百分比 */
        audioProgress: State<number>
        /** 视频下载进度百分比 */
        videoProgress: State<number>
        /** 合并进度百分比 */
        mergeProgress: State<number>
        /** 任务状态 */
        statusState: State<TaskStatus>
    })[]> = van.state([])

    constructor() {

        this.element = this.Root()
    }

    Root() {
        const _that = this
        return Route({
            rule: 'task',
            Loader() {

                return div(
                    () => div({ class: 'list-group' },
                        _that.taskList.val.map((task) => {
                            return div({
                                class: 'list-group-item list-group-item-action user-select-none vstack gap-2',
                                style: `cursor: pointer;`,
                            },
                                div(task.title),
                                div({ class: 'text-secondary', hidden: () => task.statusState.val == 'done' }, (() => {
                                    if (task.statusState.val == 'waiting') return '等待下载'
                                    if (task.statusState.val == 'error') return '下载失败'
                                    if (task.videoProgress.val == 0) {
                                        return '正在下载音频'
                                    } else if (task.mergeProgress.val == 0) {
                                        return '正在下载视频'
                                    } else if (task.statusState.val == 'running') {
                                        return '正在合并音视频'
                                    }
                                })()),
                                div({ class: 'progress', style: `height: 5px;`, hidden: () => task.statusState.val == 'done' },
                                    div({
                                        class: 'progress-bar progress-bar-striped progress-bar-animated',
                                        style: () => `width: ${(() => {
                                            if (task.videoProgress.val == 0) return task.audioProgress.val * 100
                                            if (task.mergeProgress.val == 0) return task.videoProgress.val * 100
                                            return task.mergeProgress.val * 100
                                        })()}%;`
                                    }),
                                )
                            )
                        })
                    )
                )
            },
            async onFirst() {
                if (!await checkLogin()) return
            },
            async onLoad() {
                if (!GLOBAL_HAS_LOGIN.val) return goto('login')

                const refresh = () => {
                    getActiveTask().then(activeTaskList => {
                        if (!activeTaskList) return
                        _that.taskList.val.forEach(taskInDB => {
                            activeTaskList.forEach(task => {
                                if (taskInDB.id == task.id) {
                                    taskInDB.audioProgress.val = task.audioProgress
                                    taskInDB.videoProgress.val = task.videoProgress
                                    taskInDB.mergeProgress.val = task.mergeProgress
                                    taskInDB.statusState.val = task.status
                                }
                            })
                        })
                        if (activeTaskList.filter(task => task.status == 'running').length == 0) {
                            clearInterval(timer)
                            clearInterval(halper)
                            return
                        }
                    })
                }
                refresh()
                let timer = setInterval(() => {
                    refresh()
                }, 1000)
                let halper = setInterval(() => {
                    if (now.val.split('/')[0] != 'task') {
                        clearInterval(halper)
                        clearInterval(timer)
                    }
                })
                getTaskList(0, 360).then(taskList => {
                    if (!taskList) return
                    _that.taskList.val = taskList.map(task => ({
                        ...task,
                        audioProgress: van.state(1),
                        videoProgress: van.state(1),
                        mergeProgress: van.state(1),
                        statusState: van.state('done'),
                    }))
                })
            },
        })
    }
}

export default () => new TaskRoute().element