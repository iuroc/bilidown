/// <reference types="vite/client" />

import type { workAPI } from '../../app/src/preloads/work'
import type { downloadAPI } from '../../app/src/preloads/download'
import type { settingsAPI } from '../../app/src/preloads/settings'

declare global {
    interface Window {
        workAPI: typeof workAPI
        downloadAPI: typeof downloadAPI
        settingsAPI: typeof settingsAPI
    }
}

export { }