import van from 'vanjs-core'
import { Route, goto, nowHash } from 'vanjs-router'
import { checkLogin, hasLogin } from '../mixin'
import { getQRInfo, getQRStatus } from './data'

const { div, img } = van.tags

export default () => {
    const qrSrc = van.state('')

    return Route({
        rule: 'login',
        Loader() {
            return div(
                div({ class: 'card card-body rounded-4' },
                    div({ class: 'row' },
                        div({ class: 'col-xl-3 col-lg-4 col-md-5 col-sm-6' },
                            div({ class: 'ratio ratio-1x1' },
                                img({
                                    src: qrSrc,
                                    class: 'w-100',
                                    ondragstart: event => event.preventDefault(),
                                })
                            )
                        ),
                        div({ class: 'col-xl-9 col-lg-8 col-md-7 col-sm-6' },
                            div({ class: 'vstack gap-3 h-100 px-3 justify-content-center align-items-center align-items-sm-start pb-4 pb-sm-0' },
                                div({ class: 'fs-1' }, '扫码登录'),
                                div({ class: 'fs-4' }, '使用哔哩哔哩 APP 扫码登录'),
                            )
                        )
                    )
                )
            )
        },
        async onLoad() {
            if (hasLogin.val) return goto('work')
            let qrKey = ''

            const refreshQR = async () => {
                try {
                    const qrInfo = await getQRInfo()
                    qrSrc.val = qrInfo.image
                    qrKey = qrInfo.key
                } catch (error) {
                    alert((error as Error).message)
                    clearTimeout(timer)
                }
            }

            const checkQRStatus = async () => {
                try {
                    if (await getQRStatus(qrKey)) {
                        clearTimeout(statusTimer)
                        hasLogin.val = true
                        goto('work')
                    }
                } catch (error) {
                    alert((error as Error).message)
                    clearTimeout(statusTimer)
                }
            }

            refreshQR()

            const timer = setInterval(async () => {
                if (nowHash().split('/')[0] != 'login')
                    return clearInterval(timer)
                refreshQR()
            }, 120000)

            const statusTimer = setInterval(async () => {
                if (nowHash().split('/')[0] != 'login')
                    return clearInterval(statusTimer)
                checkQRStatus()
            }, 1000)
        },
        async onFirst() {
            if (await checkLogin()) return goto('work')
        }
    })
}