import van from 'vanjs-core'

const { div } = van.tags


document.title += `Bilidown ${await window.settingsAPI.getAppVersion()} - 软件设置`

const SettingsdWindow = () => {

    return div('软件设置页面')
}

van.add(document.body, SettingsdWindow())