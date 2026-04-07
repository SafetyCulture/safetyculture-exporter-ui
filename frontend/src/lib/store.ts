import { writable } from 'svelte/store';

// shadowConfig will contain YAML configuration shadowed
let storedShadowConfig: Record<string, unknown>
try {
    storedShadowConfig = JSON.parse(localStorage.getItem("cfg") ?? '{}')
} catch (e) {
    storedShadowConfig = {}
}
export const shadowConfig = writable<Record<string, unknown>>(storedShadowConfig);
shadowConfig.subscribe((value) => {
    if (value === null) {
        value = {}
    }
    localStorage.setItem("cfg", JSON.stringify(value));
});

// template store
let storedTemplateCache: Array<Record<string, unknown>>
try {
    storedTemplateCache = JSON.parse(localStorage.getItem("templates") ?? '[]')
} catch (e) {
    storedTemplateCache = []
}
export const templateCache = writable<Array<Record<string, unknown>>>(storedTemplateCache)
templateCache.subscribe((value) => {
    if (value === null) {
        value = []
    }
    localStorage.setItem("templates", JSON.stringify(value))
})

// latest version
let storedLatestVersion: Record<string, string>
try {
    storedLatestVersion = JSON.parse(localStorage.getItem("public-version") ?? '{}')
} catch (e) {
    storedLatestVersion = {}
}
export const latestVersion = writable<Record<string, string>>(storedLatestVersion)
latestVersion.subscribe((value) => {
    if (value === null) {
        value = {}
    }
    localStorage.setItem("public-version", JSON.stringify(value))
})

// export status store
let storedExportConfig: Record<string, unknown>
try {
    storedExportConfig = JSON.parse(localStorage.getItem("export-config") ?? '{}')
} catch (e) {
    storedExportConfig = {}
}
export const exportConfig = writable<Record<string, unknown>>(storedExportConfig)
exportConfig.subscribe((value) => {
    if (value === null) {
        value = {}
    }
    localStorage.setItem("export-config", JSON.stringify(value))
})

export function emptyStores(): void {
    templateCache.set([])
    shadowConfig.set({})
    latestVersion.set({})
    exportConfig.set({})
}
