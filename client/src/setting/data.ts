import { Fields } from '.'
import { ResJSON } from '../mixin'

export const getFields = async (): Promise<Fields> => {
    const res = await fetch('/api/getFields').then(res => res.json()) as ResJSON<Fields>
    if (!res.success) throw new Error(res.message)
    return res.data
}

export const saveFields = async (fields: [string, string][]) => {
    const res = await fetch('/api/saveFields', {
        method: 'POST',
        body: JSON.stringify(fields),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => res.json()) as ResJSON
    if (!res.success) throw new Error(res.message)
    return res.message
}