import van from 'vanjs-core'

const { div, button, input } = van.tags

export default () => {

    return div({ class: 'hstack gap-3' },
        input({ placeholder: '请输入待解析的链接或代码', class: 'form-control' }),
        button({ class: 'btn btn-success text-nowrap' }, '开始解析')
    )
}