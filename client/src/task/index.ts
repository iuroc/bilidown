import van, { State } from 'vanjs-core'
import { Route, goto, now } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN, GLOBAL_HIDE_PAGE, ResJSON, VanComponent } from '../mixin'
import { getActiveTask, getTaskList, showFile } from './data'
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
                            const filename = `${task.title} ${btoa(task.id.toString()).replace(/=/g, '')}.mp4`
                            return div({
                                class: 'list-group-item vstack gap-2',
                                style: `cursor: pointer;`,
                                onclick() {
                                    if (task.statusState.val != 'done') return
                                    showFile(`${task.folder}\\${filename}`)
                                }
                            },
                                div({ class: 'small' }, filename),
                                div({ class: 'text-secondary small' },
                                    () => {
                                        if (task.statusState.val == 'waiting') return '等待下载'
                                        if (task.statusState.val == 'error') return '下载失败'
                                        if (task.videoProgress.val == 0) {
                                            return '正在下载音频'
                                        } else if (task.mergeProgress.val == 0) {
                                            return '正在下载视频'
                                        } else if (task.statusState.val == 'running') {
                                            return '正在合并音视频'
                                        } else {
                                            return task.folder
                                        }
                                    }
                                ),
                                div({
                                    class: `progress`,
                                    style: `height: 5px`,
                                    hidden: () => task.statusState.val == 'done'
                                },
                                    div({
                                        class: () => `progress-bar progress-bar-striped progress-bar-animated bg-${(() => {
                                            if (task.videoProgress.val == 0) return 'primary'
                                            if (task.mergeProgress.val == 0) return 'success'
                                            else return 'info'
                                        })()}`,
                                        style: () => {
                                            let width = 0
                                            if (task.videoProgress.val == 0) width = task.audioProgress.val * 100
                                            else if (task.mergeProgress.val == 0) width = task.videoProgress.val * 100
                                            else width = task.mergeProgress.val * 100
                                            return `width: ${width}%`
                                        }
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
                this.element.style.display = 'none'

                getTaskList(0, 360).then(taskList => {
                    if (!taskList) return
                    _that.taskList.val = taskList.map(task => ({
                        ...task,
                        audioProgress: van.state(1),
                        videoProgress: van.state(1),
                        mergeProgress: van.state(1),
                        statusState: van.state('done'),
                    }))

                    const refresh = async () => {
                        const activeTaskList = await getActiveTask()
                        if (!activeTaskList) return false
                        this.element.style.display = 'block'
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
                        }
                        return true
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
                })
            },
        })
    }
}

export default () => new TaskRoute().element