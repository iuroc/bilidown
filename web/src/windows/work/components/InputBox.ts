import van from 'vanjs-core'

const { div, button, input } = van.tags

export default () => {

    return div({ class: 'hstack gap-3' },
        input({ placeholder: '请输入待解析的链接或指令', class: 'form-control border-2', autofocus: true }),
        button({ class: 'btn btn-success text-nowrap' }, '开始解析')
    )
}