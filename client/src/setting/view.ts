import van from 'vanjs-core'
import { showErrorPage } from '../mixin'
import { SettingRoute } from '.'

const { a, button, div, input } = van.tags

export const SaveFolderSetting = (route: SettingRoute) => {
    const saveFolder = route.fields.download_folder
    const folderPickerDisabled = van.state(false)
    const buttonText = van.derive(() => folderPickerDisabled.val ? '请在弹窗中选择' : '选择目录')

    return div({ class: 'input-group' },
        div({ class: 'input-group-text' }, '下载目录'),
        input({
            class: 'form-control',
            value: saveFolder,
            oninput: event => saveFolder.val = event.target.value,
            disabled: true
        }),
        button({
            class: 'btn btn-success', onclick() {
                folderPickerDisabled.val = true
                fetch(`/api/folderPicker`).then(res => res.json()).then(res => {
                    if (!res.success) throw new Error(res.message)
                    saveFolder.val = res.data
                }).catch(error => {
                    const { message } = error as Error
                    if (message != 'Cancelled') showErrorPage(message)
                }).finally(() => {
                    folderPickerDisabled.val = false
                })
            }, disabled: folderPickerDisabled
        }, buttonText)
    )
}