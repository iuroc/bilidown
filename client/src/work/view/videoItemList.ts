import van, { State } from 'vanjs-core'
import { VideoParseResult, VideoInfoCardMode, PageInParseResult, SectionItem } from '../type'
import { VanComponent } from '../../mixin'
import { WorkRoute } from '..'

const { button, div, span } = van.tags

class VideoItemListComp implements VanComponent {
    element: HTMLElement

    constructor(
        public workRoute: WorkRoute
    ) {
        const { videoInfocardData: data } = workRoute

        this.element = div({
            hidden: () => false && data.val.pages.length <= 1,
            class: 'vstack gap-4'
        },
            div({ class: 'vstack gap-4' },
                div({ hidden: () => workRoute.allSection.val.length == 1 && workRoute.allSection.val[0].title == '正片' }, SectionTabs(this, workRoute.allSection)),
                ButtonGroup(workRoute),
                ListBox(workRoute.sectionPages),
            )
        )
    }
}

const SectionTabs = (parent: VideoItemListComp, allSection: State<SectionItem[]>) => {
    return () => div({ class: 'nav nav-underline' },
        allSection.val.map((item, index) => div({ class: 'nav-item', role: 'button' },
            div({
                tabIndex: 0,
                class: `nav-link ${parent.workRoute.sectionTabsActiveIndex.val == index ? 'active' : ''}`,
                onclick() {
                    parent.workRoute.sectionTabsActiveIndex.val = index
                },
                onkeyup(e) {
                    if (e.key == 'Enter') {
                        e.target.click()
                    }
                }
            }, () => item.title)
        ))
    )
}

const ButtonGroup = (workRoute: WorkRoute) => {
    const pages = workRoute.sectionPages
    const selectedCount = van.derive(() => pages.val.filter(page => page.selected.val).length)
    const totalCount = van.derive(() => pages.val.length)

    return div({ class: 'hstack gap-3' },
        button({
            class: 'btn btn-secondary',
            onclick() {
                pages.val.forEach(page => page.selected.val = selectedCount.val < totalCount.val)
            }
        }, () => `${selectedCount.val < totalCount.val ? '全选' : '取消全选'} (${selectedCount.val}/${totalCount.val})`),
        button({
            class: 'btn btn-primary',
            disabled: () => selectedCount.val <= 0,
            async onclick() {
                workRoute.parseModal.show()
            }
        }, '解析选中项目')
    )
}

const ListBox = (pages: State<PageInParseResult[]>) => {
    return () => div({ class: 'row gy-3 gx-3' },
        pages.val.map(page => {
            const bandgeNotNum = !page.bandge.match(/^\d+$/)
            const active = page.selected
            return div({ class: 'col-xxl-3 col-lg-4 col-md-6' },
                div({
                    tabIndex: 0,
                    class: () => `${bandgeNotNum
                        ? `vstack gap-2 justify-content-center`
                        : `hstack gap-3`
                        } shadow-sm h-100 text-break user-select-none card card-body video-item-btn bg-success bg-opacity-10 ${active.val ? 'active' : ''}`,
                    onclick() {
                        active.val = !active.val
                    },
                    onkeyup(e) {
                        if (e.key == 'Enter') {
                            active.val = !active.val
                        }
                    }
                },
                    span({ class: 'badge text-bg-success bg-opacity-75 border', hidden: bandgeNotNum }, page.bandge),
                    div(page.part),
                    div({ class: `${page.part ? 'small text-muted' : ''}`, hidden: !bandgeNotNum }, page.bandge),
                )
            )
        }),
    )
}

export default (
    workRoute: WorkRoute
) => new VideoItemListComp(workRoute).element