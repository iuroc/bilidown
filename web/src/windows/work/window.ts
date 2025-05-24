import van from 'vanjs-core'

const { div, button } = van.tags

document.title += `Bilidown 哔哩哔哩视频解析下载工具 ${await window.workAPI.getAppVersion()}`

const WorkWindow = () => {

    return div(
        button({
            onclick() {
                window.workAPI.openDownloadWindow()
            }
        }, '打开下载管理'),
        button({
            onclick() {
                window.workAPI.openSettingsWindow()
            }
        }, '打开软件设置'),
    )
}

van.add(document.body, WorkWindow())