import van from 'vanjs-core'

const { div, button } = van.tags

const DownlaodWindow = () => {

    return div(
        button({
            onclick() {

            }
        }, '打开下载管理'),
        button({
            onclick() {

            }
        }, '打开软件设置'),
    )
}

van.add(document.body, DownlaodWindow())