import van from 'vanjs-core'
import 'bootstrap/dist/css/bootstrap.min.css'

const { div } = van.tags

document.title += `下载管理`

const DownlaodWindow = () => {

    return div('下载管理页面')
}

van.add(document.body, DownlaodWindow())