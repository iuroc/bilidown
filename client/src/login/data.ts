import { ResJSON } from '../mixin'

export const getQRInfo = async () => {
    const res = await fetch('/api/getQRInfo').then(res => res.json()) as ResJSON<{ image: string, key: string }>
    if (!res.success) throw new Error(res.message)
    return res.data
}

export const getQRStatus = async (key: string) => {
    const res = await fetch('/api/getQRStatus?key=' + key).then(res => res.json()) as ResJSON
    return res.success
}