import van from 'vanjs-core'
import { Route, goto } from 'vanjs-router'
import { checkLogin, hasLogin } from '../mixin'

const { a, button, div, input } = van.tags

export default () => {
    const folderPickerDisabled = van.state(false)
    const saveFolder = van.state('')

    const SaveFolderSetting = () => div({ class: 'input-group' },
        div({ class: 'input-group-text' }, '下载目录'),
        input({
            class: 'form-control',
            value: saveFolder,
            oninput: event => saveFolder.val = event.target.value,
            disabled: 'showDirectoryPicker' in window
        }),
        button({
            class: 'btn btn-success', onclick() {
                folderPickerDisabled.val = true
                fetch(`/api/folderPicker`).then(res => res.json()).then(res => {
                    if (!res.success) throw new Error(res.message)
                    saveFolder.val = res.data
                }).catch(error => {
                    const { message } = error as Error
                    if (message != 'Cancelled') alert(message)
                }).finally(() => {
                    folderPickerDisabled.val = false
                })
            }, disabled: folderPickerDisabled
        }, '选择目录')
    )

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
        delayed: true,
        async onFirst() {
            if (!await checkLogin()) return
        },
        onLoad() {
            if (!hasLogin.val) return goto('login')
            this.show()
        },
    })
}