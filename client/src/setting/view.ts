import van from 'vanjs-core'
import { SettingRoute } from '.'
import { saveFields } from './data'

const { a, button, div, input } = van.tags

export const SaveFolderSetting = (route: SettingRoute) => {
    const saveFolder = route.fields.download_folder
    const folderPickerDisabled = van.state(false)
    const buttonText = '保存'

    return div({ class: 'input-group' },
        div({ class: 'input-group-text' }, '下载目录'),
        input({
            class: 'form-control',
            value: saveFolder,
            oninput: event => saveFolder.val = event.target.value,
        }),
        button({
            class: 'btn btn-success', onclick() {
                folderPickerDisabled.val = true
                saveFields([
                    ['download_folder', saveFolder.val]
                ]).then(message => {
                    alert(message)
                }).catch(error => {
                    if (error instanceof Error) alert(error.message)
                }).finally(() => {
                    folderPickerDisabled.val = false
                })
            }, disabled: folderPickerDisabled
        }, buttonText)
    )
}