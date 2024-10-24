import van, { State, Val } from 'vanjs-core'

const { div, img, input } = van.tags

export const VideoInfoCard = (option: {
    data: State<VideoInfoCardData>
}) => {
    const InputGroup = (title: string, value: State<string>, option?: {
        disabled?: Val<boolean>
        elementType?: 'input' | 'textarea'
    }) => {
        return div({ class: 'input-group input-group-sm' },
            div({ class: 'input-group-text' }, title),
            van.tags[option?.elementType || 'input']({
                class: 'form-control bg-white',
                disabled: option?.disabled || false,
                style: 'cursor: text;',
                value
            })
        )
    }
    const DescriptionGroup = (bottom = false) => div({
        class: () => `input-group input-group-sm flex-fill ${bottom ? 'd-lg-none' : 'd-none d-lg-flex'}`,
    },
        div({ class: 'input-group-text' }, '描述'),
        div({ class: 'form-control hstack' },
            () => option.data.val.desc.match(/^(\s*|-)$/) ? '暂无描述' : option.data.val.desc
        )
    )

    const ownerFaceHide = van.state(true)

    return div({ class: 'card border-3', hidden: () => !option.data.val.title },
        div({ class: 'card-header' }, () => option.data.val.title),
        div({ class: 'card-body vstack gap-3' },
            div({ class: 'row gx-3 gy-3' },
                // 封面
                div({ class: 'col-md-5 col-xl-4' },
                    div({ class: 'position-relative' },
                        img({
                            src: () => option.data.val.pic,
                            class: 'w-100 rounded',
                            ondragstart: event => event.preventDefault(),
                            referrerPolicy: 'no-referrer',
                            onload: () => ownerFaceHide.val = false
                        }),
                        img({
                            src: () => option.data.val.owner.face,
                            hidden: ownerFaceHide,
                            referrerPolicy: 'no-referrer',
                            ondragstart: event => event.preventDefault(),
                            style: `right: 1rem; bottom: 1rem;`,
                            class: 'rounded-3 border shadow position-absolute w-25'
                        })
                    ),
                ),
                // 字段信息
                div({ class: 'col-md-7 col-xl-8 vstack gap-2' },
                    div({ class: 'row gx-2 gy-2' },
                        div({ class: 'col-xl-7 col-xxl-8' },
                            InputGroup('作者',
                                van.derive(() => option.data.val.owner.name), { disabled: true }
                            ),
                        ),
                        div({ class: 'col-xl-5 col-xxl-4' },
                            InputGroup('发布时间',
                                van.derive(() => new Date(option.data.val.pubdate * 1000).toLocaleString()), { disabled: true }
                            )
                        ),
                        div({ class: 'col-sm col-md-12 col-lg-4' },
                            InputGroup('分辨率',
                                van.derive(() => `${option.data.val.dimension.width}x${option.data.val.dimension.height}`),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col col-lg-4' },
                            InputGroup('时长',
                                van.derive(() => `${secondToTime(option.data.val.duration)}`),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col col-lg-4' },
                            InputGroup('集数',
                                van.derive(() => option.data.val.pages.length.toString()),
                                { disabled: true }
                            )
                        ),
                    ),
                    DescriptionGroup(),
                ),
            ),
            DescriptionGroup(true),
        )
    )
}

/** 将秒数转换为 `mm:ss` */
export const secondToTime = (second: number) => {
    return `${Math.floor(second / 60)}:${(second % 60).toString().padStart(2, '0')}`
}

/** 视频信息卡片数据 */
export type VideoInfoCardData = {
    title: string
    description: string
    publishData: number
    cover: string
    duration: number
    pages: {
        cid: number
        bvid: string
        page: number
        from: string
        part: string
        duration: number
        dimension: {
            width: number
            height: number
            rotate: number
        }
    }[]
    owner: {
        mid: number
        name: string
        face: string
    }
    dimension: {
        width: number
        height: number
        rotate: number
    }
}

/** 接口返回的视频信息 */
export type VideoInfo = VideoInfoCardData & {
    aid: number
    staff: null | {
        mid: number
        title: string
        name: string
        face: string
    }[]
    title: string
    desc: string
    pubdate: number
    pic: string
    duration: number
    bvid: string
    pages: {
        cid: number
        page: number
        from: string
        part: string
        duration: number
        dimension: {
            width: number
            height: number
            rotate: number
        }
    }[]
    owner: {
        mid: number
        name: string
        face: string
    }
    dimension: {
        width: number
        height: number
        rotate: number
    }
}