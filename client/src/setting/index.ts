import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN, VanComponent } from '../mixin'
import { SaveFolderSetting } from './view'
import { getFields } from './data'
import { LoadingBox } from '../view'

const { button, div } = van.tags

export type Fields = Record<keyof SettingRoute['fields'], string>

export class SettingRoute implements VanComponent {
    element: HTMLElement

    loading = van.state(true)

    fields = {
        download_folder: van.state('')
    }

    constructor() {
        this.element = this.Root()
    }

    Root() {

        const _that = this

        return Route({
            rule: 'setting',
            Loader() {
                return div(
                    () => _that.loading.val ? LoadingBox() : '',
                    () => _that.loading.val ? '' : div({ class: 'vstack gap-4' },
                        SaveFolderSetting(_that),
                        div({ class: 'hstack gap-3' },
                            button({
                                class: 'btn btn-outline-secondary', onclick() {
                                    if (!confirm('确定要关闭软件吗?')) return
                                    fetch('/api/quit').then(res => res.json()).then(res => {
                                        if (!res.success) alert(res.message)
                                        else document.write(`<h2 style="text-align: center; padding: 30px 20px;">软件已关闭</h2>`)
                                    })
                                }
                            }, '关闭软件'),
                            button({
                                class: 'btn btn-outline-danger', onclick() {
                                    if (!confirm('确定要退出登录吗?')) return
                                    fetch('/api/logout').then(res => res.json()).then(res => {
                                        if (!res.success) alert(res.message)
                                        else location.reload()
                                    })
                                }
                            }, '退出登录'),
                        )
                    )
                )
            },
            async onFirst() {
                if (!await checkLogin()) return
                _that.loading.val = true
                getFields().then(fields => {
                    for (const key in fields) {
                        _that.fields[key as keyof Fields].val = fields[key as keyof Fields]
                    }

                    setTimeout(() => {
                        _that.loading.val = false
                    }, 200)
                })
            },
            onLoad() {
                if (!GLOBAL_HAS_LOGIN.val) return goto('login')
            },
        })
    }
}

export default () => new SettingRoute().element