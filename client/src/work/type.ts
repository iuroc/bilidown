import { State } from 'vanjs-core'

/**
 * 参考：[qn 视频清晰度标识](https://socialsisteryi.github.io/bilibili-API-collect/docs/video/videostream_url.html#qn%E8%A7%86%E9%A2%91%E6%B8%85%E6%99%B0%E5%BA%A6%E6%A0%87%E8%AF%86)
 * 
 * | 值  | 含义        |
 * | --- | ----------: |
 * | 6   | 240P 极速   |
 * | 16  | 360P 流畅   |
 * | 32  | 480P 清晰   |
 * | 64  | 720P 高清   |
 * | 74  | 720P60 高帧率 |
 * | 80  | 1080P 高清  |
 * | 112 | 1080P+ 高码率 |
 * | 116 | 1080P60 高帧率 |
 * | 120 | 4K 超清     |
 * | 125 | HDR 真彩色  |
 * | 126 | 杜比视界     |
 * | 127 | 8K 超高清    |
 */
export type VideoFormat = 6 | 16 | 32 | 64 | 74 | 80 | 112 | 116 | 120 | 125 | 126 | 127

/** 视频解析结果，数据用于 DOM 渲染，同时兼容 BV、EP、SS */
export type VideoParseResult = {
    /** 合集标题 */
    title: string
    /** 合集描述 */
    description: string
    /** 合集发布时间 */
    publishData: string
    /** 合集封面 */
    cover: string
    /** 合集总时长 */
    duration: number
    /** 分集列表 */
    pages: PageInParseResult[]

    section: SectionItem[]

    /** 合集作者 */
    owner: {
        mid: number
        name: string
        face: string
    }

    /** 合集分辨率 */
    dimension: {
        width: number
        height: number
        rotate: number
    }
    staff: string[]
    /** 更新状态 */
    status: string
    /** 地区信息 */
    areas: string[]
    /** 分类标签 */
    styles: string[]
    /** 播放页面 */
    targetURL: string
}

export type SectionItem = {
    title: string
    pages: PageInParseResult[]
}

export type PageInParseResult = {
    /** 分集 CID */
    cid: number
    /** 分集 BVID */
    bvid: string
    /** 分集在合集中的序号，从 1 开始 */
    page: number
    /** 分集标题 */
    part: string
    /** 分集时长 */
    duration: number
    /** 分集分辨率 */
    dimension: {
        width: number
        height: number
        rotate: number
    }
    /** 前置胶囊标签内容 */
    bandge: string
    /** 是否选中 */
    selected: State<boolean>
}

export type StaffItem = {
    mid: number
    title: string
    name: string
    face: string
}

type PageInVideoInfo = {
    cid: number
    page: number
    from: string
    part: string
    duration: number
    dimension: {
        width: number
        height: number
        rotate: number
    }
}

/** 接口返回的视频信息 */
export type VideoInfo = {
    aid: number
    staff: null | StaffItem[]
    title: string
    desc: string
    pubdate: number
    pic: string
    duration: number
    bvid: string
    pages: PageInVideoInfo[]
    owner: {
        mid: number
        name: string
        face: string
    }
    dimension: {
        width: number
        height: number
        rotate: number
    },
    ugc_season: {
        sections: {
            title: string
            episodes: {
                title: string
                pages: PageInVideoInfo[]
                bvid: string
            }[]
        }[] | null
        title: string
    }
}

/** 接口返回的剧集信息 */
export type SeasonInfo = {
    actors: string
    areas: {
        id: number
        name: string
    }[]
    cover: string
    evaluate: string
    publish: {
        is_finish: number
        pub_time: string
    };
    season_id: number
    season_title: string
    stat: {
        coins: number
        danmakus: number
        favorite: number
        favorites: number
        likes: number
        reply: number
        share: number
        views: number
    }
    styles: string[]
    title: string
    total: number
    episodes: Episode[]
    new_ep: {
        desc: string
        is_new: number
    }
    section: {
        title: string
        episodes: Episode[]
    }[] | null
}

export type Episode = {
    aid: number
    bvid: string
    cid: number
    cover: string
    dimension: {
        width: number
        height: number
        rotate: number
    }
    duration: number
    ep_id: number
    long_title: string
    pub_time: number
    title: string
}

export type VideoInfoCardMode = State<"video" | "season" | "hide">

export type PlayInfo = {
    accept_description: string[]
    accept_quality: number[]
    support_formats: {
        quality: number
        format: string
        new_description: string
        codecs: string[]
    }[]
    dash: {
        duration: number
        video: {
            id: number
            baseUrl: string
            backupUrl: string[]
            bandwidth: number
            mimeType: string
            codecs: string
            width: number
            height: number
            frameRate: string
            codecid: number
        }[]
        audio: {
            id: number
            baseUrl: string
            backupUrl: string[]
            bandwidth: number
            mimeType: string
            codecs: string
            width: number
            height: number
            frameRate: string
            codecid: number
        }[]
    }
}

/** 创建任务时的初始数据 */
export type TaskInitData = {
    bvid: string
    cid: number
    format: number
    title: string
    owner: string
    cover: string
}

/** 任务数据库中的数据 */
export type TaskInDB = TaskInitData & {
    id: number
    folder: string
    createAt: string
    status: TaskStatus
}

export type TaskStatus = 'done' | 'waiting' | 'running' | 'error'