import {writable} from "svelte/store";

export function trim(value, size) {
    if (value.length > size) {
        return value.substring(0, size).concat(" ...")
    }
    return value
}

let feedList = ['inspections', 'schedules', 'templates', 'template_permissions', 'sites', 'site_members', 'groups', 'group_users', 'schedule_assignees', 'schedule_occurrences', 'actions', 'action_assignees', 'issues', 'assets', 'users'];

export let feedsToExport = writable(feedList);
export let numberOfExportsCompleted = writable(0);
export let exportCompleted = writable(false);

export function markExportCompleted() {
    numberOfExportsCompleted.update((currentValue) => {
        let newValue = currentValue + 1;
        if(newValue === feedsToExport.length) {
            exportCompleted.set(true);
        }
        return newValue;
    })
}

export function resetFeedsToExport() {
    feedsToExport.set(feedList)
}