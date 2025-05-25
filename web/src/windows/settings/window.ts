import van from 'vanjs-core'
import 'bootstrap/dist/css/bootstrap.min.css'

const { div } = van.tags

document.title += `软件设置`

const SettingsdWindow = () => {

    return div('软件设置页面')
}

van.add(document.body, SettingsdWindow())