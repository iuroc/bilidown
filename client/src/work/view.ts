import van, { State, Val } from 'vanjs-core'
import { VideoInfoCardData } from './type'

const { a, button, div, img, input, span } = van.tags

export const VideoInfoCard = (option: {
    data: State<VideoInfoCardData>
    mode: State<'video' | 'season'>
}) => {
    const InputGroup = (title: Val<string>, value: State<string>, option?: {
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
    const DescriptionGroup = (bottom = false) => {
        const size = van.derive(() => option.mode.val == 'video' ? 'lg' : 'md')
        return div({
            class: () => `input-group input-group-sm flex-fill ${bottom ? `d-${size.val}-none` : `d-none d-${size.val}-flex`}`,
        },
            div({ class: 'input-group-text align-items-start' }, '描述'),
            div({ class: 'form-control' },
                () => option.data.val.description.match(/^(\s*|.)$/) ? '暂无描述' : option.data.val.description
            )
        )
    }

    const ownerFaceHide = van.state(true)

    return div({ class: 'card border-3', hidden: () => !option.data.val.title },
        div({ class: 'card-header' },
            a({
                class: 'link-dark text-decoration-none fw-bold', href: () => option.data.val.targetURL,
                target: '_blank',
            },
                () => option.data.val.title,
            )
        ),
        div({ class: 'card-body vstack gap-3' },
            div({ class: 'row gx-3 gy-3' },
                // 封面
                div({
                    class: () => option.mode.val == 'video'
                        ? 'col-md-5 col-xl-4'
                        : 'col-8 col-sm-6 mx-auto col-md-5 col-lg-3 col-xl-2'
                },
                    div({ class: 'position-relative' },
                        a({
                            href: () => option.data.val.targetURL,
                            title: () => `打开视频播放页面`,
                            target: '_blank',
                        },
                            img({
                                src: () => option.data.val.cover,
                                class: 'w-100 rounded',
                                ondragstart: event => event.preventDefault(),
                                referrerPolicy: 'no-referrer',
                                onload: () => ownerFaceHide.val = false
                            })
                        ),
                        a({
                            href: () => `https://space.bilibili.com/${option.data.val.owner.mid}`,
                            title: () => `查看用户主页：${option.data.val.owner.name}`,
                            target: '_blank',
                        },
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
                ),
                // 字段信息
                div({
                    class: () => option.mode.val == 'video'
                        ? 'col-md-7 col-xl-8 vstack gap-2'
                        : 'col-md-7 col-lg-9 col-xl-10 vstack gap-2'
                },
                    div({ class: 'row gx-2 gy-2' },
                        div({ class: 'col-xl-7 col-xxl-8' },
                            InputGroup(
                                van.derive(() => option.mode.val == 'video'
                                    ? (option.data.val.staff.length > 0 ? '制作信息' : '发布者')
                                    : '参演人员'),
                                van.derive(() => {
                                    if (option.data.val.staff.length > 0)
                                        return option.data.val.staff.map(i => i.trim()).join(', ')
                                    return option.data.val.owner.name
                                }), { disabled: true }
                            ),
                        ),
                        div({ class: 'col-xl-5 col-xxl-4' },
                            InputGroup('发布时间',
                                van.derive(() => option.data.val.publishData), { disabled: true }
                            )
                        ),
                        div({ class: 'col-sm col-md-12 col-lg-4', hidden: () => option.mode.val != 'video' },
                            InputGroup('分辨率',
                                van.derive(() => `${option.data.val.dimension.width}x${option.data.val.dimension.height}`),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col col-lg-4', hidden: () => option.mode.val != 'video' },
                            InputGroup('时长',
                                van.derive(() => `${secondToTime(option.data.val.duration)}`),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col col-lg-4', hidden: () => option.mode.val != 'video' },
                            InputGroup('集数',
                                van.derive(() => option.data.val.pages.length.toString()),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col-md-12 col-lg-4', hidden: () => option.mode.val != 'season' },
                            InputGroup('状态',
                                van.derive(() => option.data.val.status),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col-sm col-lg-4', hidden: () => option.mode.val != 'season' },
                            InputGroup('地区',
                                van.derive(() => option.data.val.areas.map(i => i.trim()).join(', ')),
                                { disabled: true }
                            )
                        ),
                        div({ class: 'col-sm col-lg-4', hidden: () => option.mode.val != 'season' },
                            InputGroup('标签',
                                van.derive(() => option.data.val.styles.join(', ')),
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
