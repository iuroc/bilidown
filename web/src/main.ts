/// <reference types="vite/client" />
import van from 'vanjs-core'
import { Route } from 'vanjs-router'
import HomeRoute from './routes/home'

const { div, button } = van.tags

const App = () => {

    return div(
        HomeRoute()
    )
}

van.add(document.body, App())