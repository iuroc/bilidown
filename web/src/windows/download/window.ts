import van from 'vanjs-core'

const { div } = van.tags

document.title += `下载管理`

const DownlaodWindow = () => {

    return div('下载管理页面')
}

van.add(document.body, DownlaodWindow())