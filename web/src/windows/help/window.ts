import van from 'vanjs-core'

const { div } = van.tags

const HelpWindow = () => {

    return div({ class: 'container-fluid p-4' },
    )
}

van.add(document.body, HelpWindow())