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