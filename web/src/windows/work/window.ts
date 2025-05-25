import van from 'vanjs-core'
import 'bootstrap/dist/css/bootstrap.min.css'
import InputBox from './components/InputBox'

const { div, button } = van.tags

document.title += `Bilidown 哔哩哔哩视频解析下载工具 ${await window.workAPI.getAppVersion()}`

const WorkWindow = () => {

    return div({ class: 'container-xxl p-3' },
        InputBox()
    )
}

van.add(document.body, WorkWindow())