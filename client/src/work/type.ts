import { State } from 'vanjs-core'

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
    pages: {
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
    }[]
    owner: {
        mid: number
        name: string
        face: string
    }
    dimension: {
        width: number
        height: number
        rotate: number
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
    }[]
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