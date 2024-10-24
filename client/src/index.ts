/// <reference types="vite/client" />
import van from 'vanjs-core'
import Header from './header'
import Work from './work'
import Task from './task'
import Login from './login'
import { redirect } from 'vanjs-router'
import 'bootstrap/dist/css/bootstrap.min.css'

const { div } = van.tags

redirect('home', 'work')

van.add(document.body,
    div({ class: 'container py-4 vstack gap-4' },
        Header(),
        Work(),
        Task(),
        Login()
    )
)