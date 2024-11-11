import { VanComponent } from '../mixin'
import van from 'vanjs-core'
import { Modal } from 'bootstrap'
import Plyr from 'plyr'
import 'plyr/dist/plyr.css'

const { a, button, div, video } = van.tags

export class PlayerComp implements VanComponent {
    element: HTMLElement
    /** 使用该对象前，请先将元素加入文档树后调用 `initPlayer` 方法 */
    player!: Plyr
    src = van.state('')
    filename = van.state('')

    constructor() {
        this.element = this.Root()
    }

    Root() {
        return video({ src: this.src })
    }

    initPlayer() {
        if (!this.element.parentNode) throw new Error('请将播放器元素加入 DOM 树')
        this.player = new Plyr(this.element, {
        })
    }
}

export class PlayerModalComp implements VanComponent {
    element: HTMLElement
    playerComp: PlayerComp
    modal: Modal

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
                this.playerComp.src.val = ''
            }
        })
    }

    Root() {
        const that = this
        return div({ class: `modal fade`, tabIndex: -1 },
            div({ class: 'modal-dialog modal-xl modal-fullscreen-xl-down' },
                div({ class: `modal-content` },
                    div({ class: `modal-header` },
                        div({ class: `h5 modal-title` }, '视频播放器'),
                        button({ class: `btn-close`, 'data-bs-dismiss': `modal` })
                    ),
                    div({ class: 'modal-body p-0' },
                        this.playerComp.element
                    ),
                    div({ class: 'modal-footer' },
                        button({ class: 'btn btn-secondary', 'data-bs-dismiss': `modal` }, '关闭'),
                        button({
                            class: 'btn btn-primary', onclick() {
                                const link = a({
                                    download: () => that.playerComp.filename.val + '.mp4',
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