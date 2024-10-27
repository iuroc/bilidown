import van, { State } from 'vanjs-core'
import { VideoParseResult, VideoInfoCardMode } from '../type'
import { VanComponent } from '../../mixin'

const { div, span } = van.tags

class VideoItemListComp implements VanComponent {
    element: HTMLElement
    constructor(
        public data: State<VideoParseResult>,
        public mode: VideoInfoCardMode
    ) {
        this.element = div({ hidden: () => false && data.val.pages.length <= 1, class: 'vstack gap-3' },
            () => div({ class: 'row gy-3 gx-3' },
                data.val.pages.map(page => {
                    const active = van.state(false)
                    return div({ class: 'col-xxl-3 col-lg-4 col-md-6' },
                        div({
                            class: () => `shadow-sm h-100 hstack gap-3 user-select-none card card-body video-item-btn bg-success bg-opacity-10 ${active.val ? 'active' : ''}`,
                            onclick() {
                                active.val = !active.val

                                console.log(page)
                            }
                        },
                            span({ class: 'badge text-bg-success bg-opacity-75 border' }, page.bandge),
                            div(page.part),
                        )
                    )
                })
            )
        )
    }
}

export default (
    data: State<VideoParseResult>,
    mode: VideoInfoCardMode
) => new VideoItemListComp(data, mode).element