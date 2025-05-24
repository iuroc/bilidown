import van from 'vanjs-core'

const { div } = van.tags

document.title += `Bilidown ${await window.downloadAPI.getAppVersion()} - 下载管理`

const DownlaodWindow = () => {

    return div('下载管理页面')
}

van.add(document.body, DownlaodWindow())