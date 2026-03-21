import van from 'vanjs-core'

const { div } = van.tags

export const LoadingBox = (color = 'primary') => div({
    class: 'py-4 hstack justify-content-center',
},
    div({
        class: `spinner-border text-${color}`,
    }, div({ class: 'visually-hidden' }, 'Loading...'))
)