import van from 'vanjs-core'
import { goto } from 'vanjs-router'

export type ResJSON<Data = null> = {
    success: boolean
    data: Data
    message: string
}

/** 创建请求超时控制器 */
export const timeoutController = (ms = 15000) => {
    const controller = new AbortController()
    const timer = setTimeout(() => {
        controller.abort(new Error('请求超时'))
    }, ms)
    return { signal: controller.signal, timer }
}

export const hasLogin = van.state(false)

export const checkLogin = async () => {
    if (hasLogin.val) return true
    const res = await fetch('/api/checkLogin').then(res => res.json()) as ResJSON
    hasLogin.val = res.success
    if (!res.success) return goto('login'), false
    return true
}