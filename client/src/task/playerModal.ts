import { VanComponent } from '../mixin'
import van from 'vanjs-core'
import { Modal } from 'bootstrap'
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'

const { a, button, div, video, audio } = van.tags

export class PlayerComp implements VanComponent {
    element: HTMLElement
    /** 使用该对象前，请先将元素加入文档树后调用 `initPlayer` 方法 */
    player!: Plyr
    videoElement: HTMLVideoElement
    audioElement: HTMLAudioElement
    videoPlayer!: Plyr
    audioPlayer!: Plyr
    type = van.state<'video' | 'audio'>('video')
    src = van.state('')
    filename = van.state('')

    constructor() {
        this.videoElement = video()
        this.audioElement = audio()
        this.element = this.Root()
    }

    Root(): HTMLElement {
        const that = this
        return div({ class: 'player-container' },
            div({ class: 'video-wrapper', hidden: () => that.type.val !== 'video' }, this.videoElement),
            div({ class: 'audio-wrapper', hidden: () => that.type.val !== 'audio' }, this.audioElement)
        )
    }

    initPlayer() {
        if (!this.element.parentNode) throw new Error('请将播放器元素加入 DOM 树')
        this.videoPlayer = new Plyr(this.videoElement, {})
        this.audioPlayer = new Plyr(this.audioElement, {})
        // 设置当前活动的player
        this.updateActivePlayer()
    }

    private updateActivePlayer() {
        this.player = this.type.val === 'video' ? this.videoPlayer : this.audioPlayer
    }

    setType(type: 'video' | 'audio') {
        this.type.val = type
        this.updateActivePlayer()
        // 清除非活动元素的src
        if (type === 'video') {
            this.audioElement.src = ''
        } else {
            this.videoElement.src = ''
        }
    }

    setSrc(src: string) {
        this.src.val = src
        if (this.type.val === 'video') {
            this.videoElement.src = src
            this.audioElement.src = ''
        } else {
            this.audioElement.src = src
            this.videoElement.src = ''
        }
    }
}

export class PlayerModalComp implements VanComponent {
    element: HTMLElement
    playerComp: PlayerComp
    modal: Modal
    type = van.state<'video' | 'audio'>('video')
    title = van.state('视频播放器')
    fileExtension = van.state('.mp4')

    constructor() {
        this.playerComp = new PlayerComp()
        this.element = this.Root()
        this.initModalEvent()

        van.add(document.body, this.element)
        this.playerComp.initPlayer()
        this.modal = new Modal(this.element)
    }

    initModalEvent() {
        this.element.addEventListener('shown.bs.modal', () => {
            this.playerComp.player.play()
        })

        this.element.addEventListener('hide.bs.modal', () => {
            if (!this.playerComp.player.stopped) {
                this.playerComp.player.stop()
                this.playerComp.setSrc('')
            }
        })
    }

    open(src: string, filename: string, type: 'video' | 'audio') {
        this.type.val = type
        this.title.val = type === 'video' ? '视频播放器' : '音频播放器'
        this.fileExtension.val = type === 'video' ? '.mp4' : '.m4a'
        this.playerComp.setType(type)
        this.playerComp.setSrc(src)
        this.playerComp.filename.val = filename
        this.modal.show()
    }

    Root() {
        const that = this
        return div({ class: `modal fade`, tabIndex: -1 },
            div({ class: 'modal-dialog modal-xl modal-fullscreen-xl-down' },
                div({ class: `modal-content` },
                    div({ class: `modal-header` },
                        div({ class: `h5 modal-title` }, () => this.title.val),
                        button({ class: `btn-close`, 'data-bs-dismiss': `modal` })
                    ),
                    div({ class: 'modal-body p-0' },
                        div({ class: () => that.type.val === 'video' ? 'ratio ratio-16x9' : '' }, this.playerComp.element),
                        div({ class: 'vstack p-3 gap-3' },
                            div({}, that.playerComp.filename)
                        )
                    ),
                    div({ class: 'modal-footer' },
                        button({ class: 'btn btn-secondary', 'data-bs-dismiss': `modal` }, '关闭'),
                        button({
                            class: 'btn btn-primary', onclick() {
                                const link = a({
                                    download: () => that.playerComp.filename.val + that.fileExtension.val,
                                    href: that.playerComp.src
                                })
                                that.modal.hide()
                                link.click()
                            }
                        }, '下载')
                    )
                )
            )
        )
    }
}