import { app } from 'electron'
import { clean, lt } from 'semver'
import { Octokit } from 'octokit'

interface LatestInfo {
    version: `${number}.${number}.${number}`
    message: string
}

export async function getLatestInfo() {
    const octokit = new Octokit()
    const response = await octokit.request('GET /repos/{owner}/{repo}/releases/latest', {
        owner: 'iuroc',
        repo: 'bilidown'
    })
    const newVersion = clean(response.data.tag_name)
    const currentVersion = clean(app.getVersion())
    if (!newVersion) throw new Error('github latest tag_name 格式错误')
    if (!currentVersion) throw new Error('package.json version 格式错误')
    return {
        version: newVersion,
        message: response.data.body,
        /** 是否比本地版本号更大 */
        isNewerThanLocal: lt(clean(app.getVersion()) as string, newVersion)
    } as LatestInfo
}