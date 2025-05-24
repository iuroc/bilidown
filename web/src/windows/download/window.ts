import van from 'vanjs-core'

const { div } = van.tags

const DownlaodWindow = () => {

    return div('下载管理页面')
}

van.add(document.body, DownlaodWindow())