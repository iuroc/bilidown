import van from 'vanjs-core'
import { sum } from '@/utils'

const { div } = van.tags

van.add(document.body, div('Hello World', sum(1, 2)))