import van from 'vanjs-core'
import { Route, goto, nowHash } from 'vanjs-router'
import { checkLogin, GLOBAL_HAS_LOGIN } from '../mixin'
import { getQRInfo, getQRStatus } from './data'

const { div, img } = van.tags

export default () => {
    const qrSrc = van.state('')
    const errorMessage = van.state('')
    const qrStatusMessage = van.state('')

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
                                    hidden: true,
                                    onload(this: HTMLImageElement) {
                                        this.hidden = false
                                    },
                                    ondragstart: event => event.preventDefault(),
                                })
                            )
                        ),
                        div({ class: 'col-xl-9 col-lg-8 col-md-7 col-sm-6' },
                            div({ class: 'vstack gap-3 h-100 px-3 justify-content-center align-items-center align-items-sm-start pb-4 pb-sm-0' },
                                div({ class: 'fs-1' }, '扫码登录'),
                                div({ class: 'fs-4' }, '使用哔哩哔哩 APP 扫码登录'),
                                div({ class: 'text-danger fw-bold', hidden: () => !errorMessage.val }, errorMessage),
                                div({ class: 'text-primary fw-bold', hidden: errorMessage },
                                    () => qrStatusMessage.val.replace('未扫码', '').replace('二维码已扫码未确认', '已扫码，请点击确认')
                                )
                            )
                        )
                    )
                )
            )
        },
        async onFirst() {
            if (await checkLogin()) return goto('work')
        },
        async onLoad() {
            // 检查登录标识，如果已经登录过了，则重定向到工作页
            if (GLOBAL_HAS_LOGIN.val) return goto('work')

            // 当前活动的二维码标识
            let qrKey = ''

            /** 刷新二维码 */
            const refreshQR = async () => {
                try {
                    const qrInfo = await getQRInfo()
                    // 更新二维码组件内容
                    qrSrc.val = qrInfo.image
                    // 更新当前活动的二维码标识
                    qrKey = qrInfo.key
                } catch (error) {
                    // 加载二维码时失败
                    errorMessage.val = `加载二维码失败，请刷新页面重试`
                    clearTimeout(timer)
                    clearTimeout(statusTimer)
                }
            }

            /** 检查二维码状态 */
            const checkQRStatus = async () => {
                try {
                    const status = await getQRStatus(qrKey)
                    if (status.success) {
                        clearTimeout(statusTimer)
                        GLOBAL_HAS_LOGIN.val = true
                        const url = new URL(window.location.href)
                        url.hash = ''
                        window.location.href = url.toString()
                    } else {
                        qrStatusMessage.val = status.message
                    }
                } catch (error) {
                    // 出错了，显示错误信息
                    errorMessage.val = `获取二维码状态失败，请刷新页面重试`
                    clearTimeout(statusTimer)
                    clearTimeout(timer)
                }
            }

            // 载入初始二维码
            refreshQR()

            /** 用于定时刷新二维码的定时器 */
            const timer = setInterval(async () => {
                if (nowHash().split('/')[0] != 'login')
                    return clearInterval(timer)
                refreshQR()
            }, 120000)

            /** 用于轮询检查二维码状态的定时器 */
            const statusTimer = setInterval(async () => {
                if (nowHash().split('/')[0] != 'login')
                    return clearInterval(statusTimer)
                checkQRStatus()
            }, 1000)
        },
    })
}