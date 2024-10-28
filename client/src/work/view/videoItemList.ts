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
        const { videoInfocardData: data, sectionTabsActiveIndex } = workRoute
        const allSection: State<SectionItem[]> = van.derive(() => [{ title: '正片', pages: data.val.pages }].concat(data.val.section))
        const sectionPages = van.derive(() => allSection.val[sectionTabsActiveIndex.val].pages || [])

        this.element = div({
            hidden: () => false && data.val.pages.length <= 1,
            class: 'vstack gap-4'
        },
            div({ class: 'vstack gap-4' },
                div({ hidden: () => allSection.val.length <= 1 }, SectionTabs(this, allSection)),
                ButtonGroup(sectionPages),
                ListBox(sectionPages),
            )
        )
    }
}

const SectionTabs = (parent: VideoItemListComp, allSection: State<SectionItem[]>) => {
    return () => div({ class: 'nav nav-underline' },
        allSection.val.map((item, index) => div({ class: 'nav-item user-select-none', role: 'button' },
            div({
                class: `nav-link ${parent.workRoute.sectionTabsActiveIndex.val == index ? 'active' : ''}`,
                onclick() {
                    parent.workRoute.sectionTabsActiveIndex.val = index
                }
            }, () => item.title)
        ))
    )
}

const ButtonGroup = (pages: State<PageInParseResult[]>) => {
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
                    class: () => `${bandgeNotNum
                        ? `vstack gap-2 justify-content-center`
                        : `hstack gap-3`
                        } shadow-sm h-100 user-select-none card card-body video-item-btn bg-success bg-opacity-10 ${active.val ? 'active' : ''}`,
                    onclick() {
                        active.val = !active.val
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