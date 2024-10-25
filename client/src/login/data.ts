import { ResJSON } from '../mixin'

/** 获取新的二维码信息，包含二维码 Base64 数据和二维码 Key */
export const getQRInfo = async () => {
    const res = await fetch('/api/getQRInfo').then(res => res.json()) as ResJSON<{ image: string, key: string }>
    if (!res.success) throw new Error(res.message)
    return res.data
}

/** 根据二维码 Key 获取二维码状态 */
export const getQRStatus = async (key: string) => {
    return await fetch('/api/getQRStatus?key=' + key).then(res => res.json()) as ResJSON
}