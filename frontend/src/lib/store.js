import {writable} from 'svelte/store';

// shadowConfig will contain YAML configuration shadowed
const storedShadowConfig = JSON.parse(localStorage.getItem("cfg"))
export const shadowConfig = writable(storedShadowConfig);
shadowConfig.subscribe(value => {
    localStorage.setItem("cfg", JSON.stringify(value));
});

// template store
const storedTemplateCache = JSON.parse(localStorage.getItem("templates"))
export const templateCache = writable(storedTemplateCache)
templateCache.subscribe(value => {
    localStorage.setItem("templates", JSON.stringify(value))
})
