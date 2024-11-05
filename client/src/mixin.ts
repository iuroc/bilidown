import van from 'vanjs-core'
import { goto } from 'vanjs-router'

export type ResJSON<Data = null> = {
    success: boolean
    data: Data
    message: string
}

/** 创建请求超时控制器 */
export const timeoutController = (ms = 15000): {
    signal: AbortSignal
    timer: number
} => {
    const controller = new AbortController()
    const timer = setTimeout(() => {
        controller.abort(new Error('请求超时'))
    }, ms)
    return { signal: controller.signal, timer }
}

/** 全局登录状态 */
export const GLOBAL_HAS_LOGIN = van.state(false)

/** 设置页面整体隐藏，该状态值用于设置根 DOM 的 `hidden` 属性 */
export const GLOBAL_HIDE_PAGE = van.state(true)

/** 全局错误信息，用于在统一的错误提示页面中展示 */
export const GLOBAL_ERROR_MESSAGE = van.state('')

/** 跳转到统一的错误提示页面 */
export const showErrorPage = (message: string) => {
    GLOBAL_HIDE_PAGE.val = true
    GLOBAL_ERROR_MESSAGE.val = message
    goto('error')
}

/** 检查后端是否登录，如果未登录，则跳转到登录页
 * 
 * 登录成功或失败，都将更新 `GLOBAL_HAS_LOGIN` 的值，并返回 `boolean`，并将 `GLOBAL_HIDE_PAGE` 设置为 `false`
 * 
 * 用法注意：每个路由 `onFirst` 和 `onLoad` 只能其中一个使用该函数，否则会导致执行两次请求
 * 一般在 `onFirst` 中执行本方法，在 `onLoad` 中执行 `if (!GLOBAL_HAS_LOGIN.val) return goto('login')`
 */
export const checkLogin = async (): Promise<boolean> => {
    if (GLOBAL_HAS_LOGIN.val) return GLOBAL_HIDE_PAGE.val = false, true
    const res = await fetch('/api/checkLogin').then(res => res.json()) as ResJSON
    GLOBAL_HAS_LOGIN.val = res.success
    GLOBAL_HIDE_PAGE.val = false
    if (!res.success) return goto('login'), false
    return true
}

export interface VanComponent {
    element: HTMLElement
}

export const formatSeconds = (seconds: number) => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = seconds % 60

    let str = ''
    if (hours > 0) str += `${hours}时`
    if (minutes > 0) str += `${minutes}分`
    if (secs > 0) str += `${secs}秒`

    return str
}