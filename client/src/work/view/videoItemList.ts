import van, { State } from 'vanjs-core'
import { VideoParseResult, VideoInfoCardMode } from '../type'
import { VanComponent } from '../../mixin'

const { div } = van.tags

class VideoItemListComp implements VanComponent {
    element: HTMLElement
    constructor(
        public data: State<VideoParseResult>,
        public mode: VideoInfoCardMode
    ) {
        this.element = div('分集列表')
    }
}

export default (
    data: State<VideoParseResult>,
    mode: VideoInfoCardMode
) => new VideoItemListComp(data, mode).element