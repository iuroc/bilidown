import van from 'vanjs-core'

const { div } = van.tags

const App = () => {
    return div('Hello World')
}

van.add(document.body, App())